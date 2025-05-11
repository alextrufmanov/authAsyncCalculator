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
func POSTLoginHandler(s *Storage, dbs *DBStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.POSTLoginRequestBody
		log.Printf("POST: /api/v1/login")
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				// log.Printf("  Body = %v", requestBody)
				// pswHash, res := Hash(requestBody.Login, requestBody.Password)
				// if res {
				// 	log.Printf("  pswHash = %s", pswHash)
				// 	token, res := dbs.GetToken(requestBody.Login, pswHash)
				// 	if res {
				// 		json.NewEncoder(w).Encode(models.POSTLoginAnswerBody{Token: token})
				// 		return
				// 	}
				token, res := dbs.GetToken(requestBody.Login, requestBody.Password)
				if res {
					json.NewEncoder(w).Encode(models.POSTLoginAnswerBody{Token: token})
					return
				}
			}
		}
		http.Error(w, "{\"error\":\"invalid login or pasword\"}", http.StatusUnauthorized)
	}
}
