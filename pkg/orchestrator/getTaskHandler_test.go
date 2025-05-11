package orchestrator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

func TestGETTaskHandler(t *testing.T) {

	cases := []struct {
		expressionStr string
		responseBody  models.GETTaskAnswerBody
		status        int
	}{
		{
			expressionStr: "2*2",
			responseBody:  models.GETTaskAnswerBody{Task: models.Task{Id: 2, Arg1: 2, Arg2: 2, Operation: "*", OperationTime: 5000, Status: models.TaskStatusReady}},
			status:        http.StatusOK,
		},
		{
			expressionStr: "",
			responseBody:  models.GETTaskAnswerBody{Task: models.Task{Status: models.TaskStatusWait}},
			status:        http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run("Test GETTaskHandler", func(t *testing.T) {

			cfg := config.NewCfg()
			dbs, _ := NewDBStorage()
			storage := NewStorage(*cfg, dbs)

			storage.AppendExpression(1, 1, tc.expressionStr)
			expression := storage.newExpression(2, 1, tc.expressionStr)
			task := storage.newTask(expression, "*", make(chan float64, 1))
			task.Status = tc.responseBody.Task.Status
			task.Arg1 = tc.responseBody.Task.Arg1
			task.Arg2 = tc.responseBody.Task.Arg2
			storage.tasks[task.Id] = task

			req, err := http.NewRequest("GET", "/internal/task", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			r := httptest.NewRecorder()
			handler := http.HandlerFunc(GETTaskHandler(storage))
			handler.ServeHTTP(r, req)

			status := r.Code
			if status != tc.status {
				t.Errorf("Error status %v, want %v", status, tc.status)
			}

			if status == http.StatusOK {

				var responseBody models.GETTaskAnswerBody
				err = json.NewDecoder(r.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				responseTask := responseBody.Task
				wantTask := tc.responseBody.Task

				if responseTask.Arg1 != wantTask.Arg1 {
					t.Errorf("task.Arg1 %v, want %v.", responseTask.Arg1, wantTask.Arg1)
				}

				if responseTask.Arg2 != wantTask.Arg2 {
					t.Errorf("task.Arg2 %v, want %v.", responseTask.Arg2, wantTask.Arg2)
				}

				if responseTask.Operation != wantTask.Operation {
					t.Errorf("task.Operation %v, want %v.", responseTask.Operation, wantTask.Operation)
				}

				if responseTask.OperationTime != wantTask.OperationTime {
					t.Errorf("task.OperationTime %v, want %v.", responseTask.OperationTime, wantTask.OperationTime)
				}
			}
		})
	}
}
