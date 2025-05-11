package models

// Константы возможных состояний арифметических выражений
const (
	ExpressionStatusReady     = "ready"
	ExpressionStatusCalculate = "calculate"
	ExpressionStatusFailed    = "failed"
	ExpressionStatusSuccess   = "success"
)

// Константы возможных состояний задач вычисления арифметических выражений
const (
	TaskStatusWait      = "wait"
	TaskStatusReady     = "ready"
	TaskStatusCalculate = "calculate"
	TaskStatusFailed    = "failed"
	TaskStatusSuccess   = "success"
)

type User struct {
	Id      int64
	Login   string
	PswHash string
}

// Структура арифметического выражения
type Expression struct {
	Id         int32   `json:"id"`
	UserId     int32   `json:"userId"`
	Expression string  `json:"expression"`
	Status     string  `json:"status"`
	Result     float64 `json:"result"`
	Tasks      []*Task `json:"-"`
}

// Структура задач вычисления арифметических выражений
type Task struct {
	Owner         *Expression  `json:"-"`
	Id            int32        `json:"id"`
	Arg1          float64      `json:"arg1"`
	Arg2          float64      `json:"arg2"`
	Operation     string       `json:"Operation"`
	OperationTime int32        `json:"operation_time"`
	Status        string       `json:"-"`
	Result        float64      `json:"-"`
	ResultChan    chan float64 `json:"-"`
}
