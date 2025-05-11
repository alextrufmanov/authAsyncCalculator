package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

// Функция возвращает обработчик GET запроса эндпоинта /api/v1/expressions,
// запрос у оркестратора информации обо всех арифметических выражениях пользователя
func GETExpressionsHandler(s *Storage, dbs *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.GETExpressionsRequestBody
		log.Printf("GET: /api/v1/expressions")
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				user_id, err := JWTToUserId(requestBody.Token)
				if err == nil {
					json.NewEncoder(w).Encode(models.GETexpressionsAnswerBody{Expressions: dbs.GetAllExpressions(user_id)})
					return
				}
			}
		}
		SendForbiddenError403(w)
	}
}
