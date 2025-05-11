package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
	"github.com/gorilla/mux"
)

// Функция возвращает обработчик GET запроса эндпоинта /api/v1/expressions/{id},
// запрос у оркестратора информации о арифметическом выражении с указанным id
func GETExpressionByIDHandler(s *Storage, dbs *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.Atoi(idStr)
		if err == nil {
			log.Printf("GET: /api/v1/expressions/%d", id)

			var requestBody models.GETExpressionByIDRequestBody
			bodyBytes, err := io.ReadAll(r.Body)
			if err == nil {
				if json.Unmarshal(bodyBytes, &requestBody) == nil {
					user_id, err := JWTToUserId(requestBody.Token)
					if err == nil {
						expression, r := dbs.GetExpressionByID(int32(id), user_id)
						if r {
							json.NewEncoder(w).Encode(models.GETexpressionByIDAnswerBody{Expression: expression})
							return
						} else {
							SendNotFoundError404(w)
							return
						}
					}
				}
			}

		}
		SendForbiddenError403(w)
	}
}
