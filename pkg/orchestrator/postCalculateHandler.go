package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

// Функция возвращает обработчик POST запроса эндпоинта /api/v1/calculate
// запрос на асинхронное вычисление нового арифметического выражения
func POSTCalculateHandler(s *Storage, dbs *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.POSTCalculateRequestBody
		log.Printf("POST: /api/v1/calculate")
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				// log.Printf("  Body = %v", requestBody)
				user_id, err := JWTToUserId(requestBody.Token)
				if err == nil {
					id, err := dbs.AppendExpression(user_id, requestBody.Expression)
					if err == nil {
						_, res := s.AppendExpression(id, user_id, requestBody.Expression)
						if res {
							json.NewEncoder(w).Encode(models.POSTCalculateAnswerBody{Id: id})
							return
						} else {
							dbs.UpdateExpressionStatus(id, models.ExpressionStatusFailed, 0)
						}
					}
				}
			}
		}
		SendInvalidDataError422(w)
	}
}
