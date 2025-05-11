package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

// Функция возвращает обработчик POST запроса эндпоинта /api/v1/register
// запрос на асинхронное вычисление нового арифметического выражения
func POSTRegisterHandler(s *Storage, dbs *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.POSTRegisterRequestBody
		log.Printf("POST: /api/v1/register")
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				// log.Printf("  Body = %v", requestBody)
				// pswHash, res := Hash(requestBody.Login, requestBody.Password)
				// if res {
				// 	log.Printf("  pswHash = %s", pswHash)
				// 	_, res := dbs.InsertUser(requestBody.Login, pswHash)
				// 	if res {
				// 		return
				// 	}
				// }
				_, res := dbs.InsertUser(requestBody.Login, requestBody.Password)
				if res {
					return
				}

			}
		}
		http.Error(w, "{\"error\":\"Can't create user\"}", http.StatusUnprocessableEntity)
	}
}
