package models

// Структура тела POST запроса на регистрацию нового пользователя
type POSTRegisterRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Структура тела POST запроса на логин пользователя
type POSTLoginRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Структура тела ответа на POST запрос  на логин пользователя
type POSTLoginAnswerBody struct {
	Token string `json:"token"`
}

// Структура тела POST запроса на регистрацию нового пользователя
type POSTCalculateRequestBody struct {
	Token      string `json:"token"`
	Expression string `json:"expression"`
}

// Структура тела ответа на POST запрос на вычисление арифметического выражения
type POSTCalculateAnswerBody struct {
	Id int32 `json:"id"`
}

// Структура тела GET запроса списка арифметических выражений
type GETExpressionsRequestBody struct {
	Token string `json:"token"`
}

// Структура тела ответа на GET запрос списка арифметических выражений (обработанных и находящихся в обработке)
type GETexpressionsAnswerBody struct {
	Expressions []Expression `json:"expressions"`
}

// Структура тела GET запроса арифметического выражения по его Id
type GETExpressionByIDRequestBody struct {
	Token string `json:"token"`
}

// Структура тела ответа на GET запрос арифметического выражения по его Id
type GETexpressionByIDAnswerBody struct {
	Expression Expression `json:"expression"`
}

// Структура тела ответа на внутренний GET запрос задачи агентом
type GETTaskAnswerBody struct {
	Task Task `json:"task"`
}

// Структура тела внутреннего POST запроса передачи результатов вычислений агентом
type POSTTaskResultRequestBody struct {
	Id      int32   `json:"id"`
	Result  float64 `json:"result"`
	Success bool    `json:"success"`
}
