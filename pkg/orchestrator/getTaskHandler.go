package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

// Функция возвращает обработчик GET запроса эндпоинта /internal/task,
// получение агентом у оркестратора очередной задачи на выполнение
func GETTaskHandler(s *Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("GET: /internal/task")
		task, res := s.GetTask()
		if res {
			json.NewEncoder(w).Encode(models.GETTaskAnswerBody{Task: task})
			return
		}
		SendNotFoundError404(w)
	}
}
