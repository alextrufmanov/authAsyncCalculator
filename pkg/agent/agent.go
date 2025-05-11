package agent

import (
	"context"
	"log"
	"time"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
	"github.com/alextrufmanov/asyncCalculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// // Функция пытается получить от оркестратора очередную задачу
// func httpGetTask(addr string) (models.Task, bool) {
// 	var answerBody models.GETTaskAnswerBody
// 	response, err := http.Get(fmt.Sprintf("http://%s/internal/task", addr))
// 	if err == nil {
// 		if response.StatusCode == http.StatusOK {
// 			bodyBytes, err := io.ReadAll(response.Body)
// 			if err == nil {
// 				if json.Unmarshal(bodyBytes, &answerBody) == nil {
// 					return answerBody.Task, true
// 				}
// 			}
// 		}
// 	}
// 	return models.Task{}, false
// }

// Функция пытается получить от оркестратора очередную задачу
func grpcGetTask(addr string) (models.Task, bool) {
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// log.Println("Could not connect to grpc server: ", err)
		return models.Task{}, false
	}
	defer conn.Close()
	grpcClient := proto.NewAsyncCalculatorServiceClient(conn)
	task, err := grpcClient.GetTask(context.TODO(), &proto.Undef{})
	if err != nil {
		// log.Println("failed invoking PostTaskResulp: ", err)
		return models.Task{}, false
	}
	return models.Task{
		Id:            task.Id,
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     task.Operation,
		OperationTime: task.OperationTime,
	}, true

	// area, err := grpcClient.PostTaskResulp(context.TODO(), &proto.Undef{})
}

// // Функция отправляет результат решения задачи оркестратору
// func httpPostTaskResul(addr string, id int, result float64, success bool) bool {
// 	requestBody, err := json.Marshal(models.POSTTaskResultRequestBody{Id: id, Result: result, Success: success})
// 	if err == nil {
// 		response, err := http.Post(fmt.Sprintf("http://%s/internal/task", addr), "application/json", bytes.NewBuffer(requestBody))
// 		if err == nil {
// 			return response.StatusCode == http.StatusOK
// 		}
// 	}
// 	return false
// }

// Функция отправляет результат решения задачи оркестратору
func grpcPostTaskResul(addr string, id int32, result float64, success bool) bool {
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Could not connect to grpc server: ", err)
		return false
	}
	defer conn.Close()
	grpcClient := proto.NewAsyncCalculatorServiceClient(conn)
	_, err = grpcClient.PostTaskResult(context.TODO(), &proto.TaskResult{
		Id:      id,
		Result:  result,
		Success: success,
	})
	if err != nil {
		log.Println("failed invoking PostTaskResulp: ", err)
		return false
	}
	return true
}

// Функция отправляет результат решения задачи оркестратору
func calculate(task *models.Task) bool {
	// имитируем "длительную" задачу
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	// выполняем задачу
	switch task.Operation {
	case "+":
		task.Result = task.Arg1 + task.Arg2
	case "-":
		task.Result = task.Arg1 - task.Arg2
	case "*":
		task.Result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			log.Printf("Задача %d: (%f) %s (%f) => Деление на ноль", task.Id, task.Arg1, task.Operation, task.Arg2)
			return false
		}
		task.Result = task.Arg1 / task.Arg2
	default:
		log.Printf("Задача %d: (%f) %s (%f) => Неподдерживаемый оператор", task.Id, task.Arg1, task.Operation, task.Arg2)
		return false
	}
	log.Printf("Задача %d: (%f) %s (%f) => (%f)", task.Id, task.Arg1, task.Operation, task.Arg2, task.Result)
	return true
}

// // Функция создания и запуска агентов
// func StartHttpAgents(cfg config.Cfg) {
// 	log.Printf("Agent started (%s).", cfg.Addr)
// 	// запускаем указанное количество вычислитетлей в отдельных горутинах
// 	for range cfg.ComputingPower {
// 		go func() {
// 			for {
// 				// пытаемся получить от оркестратора очередную задачу
// 				task, r := httpGetTask(cfg.Addr)
// 				if r {
// 					// если задача получена, то решаем ее
// 					success := calculate(&task)
// 					// отправляем результат решения задачи оркестратору c
// 					httpPostTaskResul(cfg.Addr, task.Id, task.Result, success)
// 				}
// 				time.Sleep(time.Duration(500) * time.Millisecond)
// 			}
// 		}()
// 	}
// 	log.Printf("%d calculators was started", cfg.ComputingPower)
// 	select {}
// }

// Функция создания и запуска агентов
func StartGrpcAgents(cfg config.Cfg) {
	log.Printf("Agent started (%s).", cfg.GrpcAddr)
	// запускаем указанное количество вычислитетлей в отдельных горутинах
	for range cfg.ComputingPower {
		go func() {
			for {
				// пытаемся получить от оркестратора очередную задачу
				task, r := grpcGetTask(cfg.GrpcAddr)
				if r {
					// если задача получена, то решаем ее
					success := calculate(&task)
					// отправляем результат решения задачи оркестратору c
					grpcPostTaskResul(cfg.GrpcAddr, task.Id, task.Result, success)
				}
				time.Sleep(time.Duration(500) * time.Millisecond)
			}
		}()
	}
	log.Printf("%d calculators was started", cfg.ComputingPower)
	select {}
}
