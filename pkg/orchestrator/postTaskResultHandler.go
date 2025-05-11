package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

// Функция возвращает обработчик POST запроса эндпоинта /internal/task,
// оповещение оркестратора агентом о результатах выполнения задачи
func POSTTaskResultHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody models.POSTTaskResultRequestBody
		log.Printf("POST: /internal/task")
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			if json.Unmarshal(bodyBytes, &requestBody) == nil {
				// log.Printf("  Body = %v", requestBody)
				if s.SetTaskResult(requestBody.Id, requestBody.Result, requestBody.Success) {
				} else {
					SendNotFoundError404(w)
				}
				return
			}
		}
		SendInternalError500(w)
	}
}
