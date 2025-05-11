package orchestrator

import (
	"strconv"
	"strings"
	"sync"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

// Структура хранилища данных (структур арифметических выражений и задач агенту)
type Storage struct {
	db          *DBStorage
	config      config.Cfg
	mutex       sync.Mutex
	expressions map[int32]*models.Expression
	tasks       map[int32]*models.Task
	genIdMutex  sync.Mutex
	genId       int32
}

// Функция создает новое хранилище данных
func NewStorage(cfg config.Cfg, dbs *DBStorage) *Storage {
	s := Storage{
		db:          dbs,
		config:      cfg,
		expressions: make(map[int32]*models.Expression),
		tasks:       make(map[int32]*models.Task),
		genId:       0,
		genIdMutex:  sync.Mutex{},
		mutex:       sync.Mutex{},
	}
	for _, expression := range dbs.GetActiveExpressions() {
		_, res := s.AppendExpression(expression.Id, expression.UserId, expression.Expression)
		if !res {
			dbs.UpdateExpressionStatus(expression.Id, models.ExpressionStatusFailed, 0)
		}
	}
	return &s
}

// Функция создает новый id
func (s *Storage) newId() int32 {
	s.genIdMutex.Lock()
	defer s.genIdMutex.Unlock()
	s.genId++
	return s.genId
}

// Функция создает (но не добавляет) в хранилище новое арифметическое выражение
func (s *Storage) newExpression(id int32, userId int32, expression string) *models.Expression {
	return &models.Expression{
		Id:         id,
		UserId:     userId,
		Expression: expression,
		// Id:         s.newId(),
		Status: models.ExpressionStatusReady,
		Result: 0,
		Tasks:  make([]*models.Task, 0),
	}
}

// Функция создает (но не добавляет) в хранилище новую задачу агента
func (s *Storage) newTask(expression *models.Expression, operation string, resultChan chan float64) *models.Task {
	var timeout int
	switch operation {
	case "+":
		timeout = s.config.AddTimeout
	case "-":
		timeout = s.config.SubTimeout
	case "*":
		timeout = s.config.MltTimeout
	case "/":
		timeout = s.config.DivTimeout
	}
	task := models.Task{
		Owner:         expression,
		Id:            s.newId(),
		Arg1:          0,
		Arg2:          0,
		Operation:     operation,
		OperationTime: int32(timeout),
		Status:        models.TaskStatusWait,
		Result:        0,
		ResultChan:    resultChan,
	}
	expression.Tasks = append(expression.Tasks, &task)
	return &task
}

// Функция добавляет в хранилище новое арифметическое выражение, создает
// связанные с этим выражением асинхронные задачи агенту
func (s *Storage) AppendExpression(id int32, userId int32, expressionStr string) (int32, bool) {
	var stack [](chan float64)

	// подготавливаем выражение к преобразованию в RPN
	items, err := Split(expressionStr)
	if err != nil {
		return -1, false
	}

	// преобразуем выражение в RPN
	rpm, err := ToRPM(items)
	if err != nil {
		return -1, false
	}

	// создаем новое выражение
	expression := s.newExpression(id, userId, expressionStr)

	// создаем задачи агенту для асинхронного вычисления преобразованного в
	// RPN арифметического выражения
	for _, item := range rpm {
		if strings.Contains(operations, item) {

			// токен RPN - оператор, создаем новую задачу,
			// в стеке уже должны быть минимум 2 канала, по которым поступят
			// результаты исполнения связанных подзадач
			if len(stack) < 2 {
				return -1, false
			}

			// получаем из стека каналы с аргументами новой задачи
			arg1Chan := stack[len(stack)-2]
			arg2Chan := stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			// создаем канал результата новой задачи и помещаем его в стек
			resultChan := make(chan float64, 1)
			stack = append(stack, resultChan)

			// создаем новую задачу для агента
			newTask := s.newTask(expression, item, resultChan)

			// запускаем в отдельном потоке ожидание результатов исполнегия
			// связанных подзадач
			go func(t *models.Task, arg1 chan float64, arg2 chan float64) {
				// ожидаем поступления результатов исполнегия связанных задач
				arg1Value := <-arg1
				arg2Value := <-arg2

				s.mutex.Lock()
				defer s.mutex.Unlock()

				// устанавливаем полученные аргументы и статус задачи "ready"
				t.Arg1 = arg1Value
				t.Arg2 = arg2Value

				if t.Status == models.TaskStatusWait {
					t.Status = models.TaskStatusReady
				}

			}(newTask, arg1Chan, arg2Chan)

		} else {
			// токен RPN - не оператор - предполагаем, что это число
			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return -1, false
			}
			// помещаем число в новый канал, а канал в стек, для чисел
			// создавать задачи не нужно
			newChan := make(chan float64, 1)
			newChan <- value
			stack = append(stack, newChan)
		}
	}

	// RPN полностью "разобрана", в стеке должен остаться только один канал,
	// в который поступит результат корневой задачи, т.е. результат всего
	// арифметического выражения
	if len(stack) != 1 {
		return -1, false
	}

	// запускаем в отдельном потоке ожидание результата исполнегия корневой
	// задачи арифметического выражения
	go func(expression *models.Expression, resultChan chan float64) {
		// ожидаем поступления результата исполнегия корневой задачи
		resultValue := <-resultChan

		s.mutex.Lock()
		defer s.mutex.Unlock()

		// устанавливаем полученный результат и статус "success" или "failed"
		if expression.Status == models.ExpressionStatusCalculate {
			expression.Result = resultValue
			expression.Status = models.ExpressionStatusSuccess
			s.db.UpdateExpressionStatus(expression.Id, models.ExpressionStatusSuccess, resultValue)
		}

		// удаляем из хранилища связанные с вычислением данного арифметического
		//  выражения задачи
		// for _, task := range expression.Tasks {
		// 	delete(s.tasks, task.Id)
		// }
	}(expression, stack[0])

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// добавляем в хранилище подготовленные арифметическое вырадение и
	// связанные с ним задачи агента
	s.expressions[expression.Id] = expression
	for _, task := range expression.Tasks {
		s.tasks[task.Id] = task
	}

	return expression.Id, true
}

// Функция возвращает из хранилища арифметическое выражение с указанным Id
func (s *Storage) GetExpressionByID(id int32) (models.Expression, bool) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	expression, r := s.expressions[id]
	if r {
		return *expression, r
	} else {
		return models.Expression{}, r
	}
}

// Функция возвращает из хранилища все арифметические выражения
func (s *Storage) GetAllExpressions() []models.Expression {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := make([]models.Expression, 0)
	for _, expression := range s.expressions {
		r = append(r, *expression)
	}
	return r
}

// Функция возвращает из хранилища первую готовую к запуску задачу агенту,
// изменяет ее статус и статус соответствующего арифметического выражения
// на "calculate"
func (s *Storage) GetTask() (models.Task, bool) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range s.tasks {
		if task.Status == models.TaskStatusReady {
			task.Status = models.TaskStatusCalculate
			task.Owner.Status = models.ExpressionStatusCalculate
			s.db.UpdateExpressionStatus(task.Owner.Id, models.ExpressionStatusCalculate, 0)
			return *task, true
		}
	}
	return models.Task{}, false
}

// Функция устанавливает результат выполнения агентом указанной задачи,
// изменяет ее статус и передает результат следующей задаче или
// арифметическому выражению, если задача является корневой
func (s *Storage) SetTaskResult(id int32, result float64, success bool) bool {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	task, res := s.tasks[id]
	if res {
		if task.Status == models.TaskStatusCalculate {
			if success {
				task.Result = result
				task.Status = models.TaskStatusSuccess
				task.ResultChan <- task.Result
			} else {
				task.Status = models.TaskStatusFailed
				task.Owner.Status = models.ExpressionStatusFailed
				s.db.UpdateExpressionStatus(task.Owner.Id, models.ExpressionStatusFailed, 0)
			}
			return true
		}
	}
	return false
}
