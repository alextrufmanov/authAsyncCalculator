package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

func TestPOSTTaskResultHandler(t *testing.T) {

	cases := []struct {
		expressionStr string
		requestBody   models.POSTTaskResultRequestBody
		status        int
	}{
		{
			expressionStr: "(2+2*2)/2",
			requestBody:   models.POSTTaskResultRequestBody{Id: 2, Result: 4, Success: true},
			status:        http.StatusOK,
		},
		{
			expressionStr: "(2+2*2)/2",
			requestBody:   models.POSTTaskResultRequestBody{Id: 2, Result: 0, Success: false},
			status:        http.StatusOK,
		},
		{
			expressionStr: "(2+2*2)/2",
			requestBody:   models.POSTTaskResultRequestBody{Id: 4, Result: 3, Success: true},
			status:        http.StatusNotFound,
		},
		{
			expressionStr: "(2+2*2)/2",
			requestBody:   models.POSTTaskResultRequestBody{Id: 1, Result: -1, Success: true},
			status:        http.StatusNotFound,
		},
	}

	config := config.NewCfg()
	dbs, _ := NewDBStorage()

	for _, tc := range cases {
		t.Run("Test POSTTaskResultHandler", func(t *testing.T) {

			POSTTaskResultHandlerStorage := NewStorage(*config, dbs)
			_, res := POSTTaskResultHandlerStorage.AppendExpression(1, 1, tc.expressionStr)
			if !res {
				t.Fatalf("failed to append expression: %v", tc.expressionStr)
			}
			POSTTaskResultHandlerStorage.tasks[2].Status = models.TaskStatusCalculate

			requestBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			r := httptest.NewRecorder()

			handler := http.HandlerFunc(POSTTaskResultHandler(POSTTaskResultHandlerStorage))
			handler.ServeHTTP(r, req)

			status := r.Code
			if status != tc.status {
				t.Errorf("Error task (%d) status %v, want %v", tc.requestBody.Id, status, tc.status)
			}

			if status == http.StatusOK {
				task, res := POSTTaskResultHandlerStorage.tasks[tc.requestBody.Id]
				if res {
					if task.Result != tc.requestBody.Result {
						t.Errorf("task.Result %v, want %v.", task.Result, tc.requestBody.Result)
					}
					if tc.requestBody.Success && task.Status != models.TaskStatusSuccess {
						t.Errorf("task.Status %v, want %v.", task.Status, models.ExpressionStatusSuccess)
					}
					if !tc.requestBody.Success && task.Status != models.TaskStatusFailed {
						t.Errorf("task.Status %v, want %v.", task.Status, models.ExpressionStatusFailed)
					}
				} else {
					t.Errorf("Task %d not exist.", tc.requestBody.Id)
				}
			}

		})
	}
}
