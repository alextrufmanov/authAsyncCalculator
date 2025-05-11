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

func TestPOSTCalculateHandler(t *testing.T) {

	cases := []struct {
		requestBody      models.POSTCalculateRequestBody
		status           int
		responseBody     models.POSTCalculateAnswerBody
		expressionsCount int
		tasksCount       int
	}{
		{
			requestBody:      models.POSTCalculateRequestBody{Expression: "2+2"},
			status:           http.StatusOK,
			responseBody:     models.POSTCalculateAnswerBody{Id: 1},
			expressionsCount: 1,
			tasksCount:       1,
		},
		{
			requestBody:      models.POSTCalculateRequestBody{Expression: "(-1+7)+(12-7)"},
			status:           http.StatusOK,
			responseBody:     models.POSTCalculateAnswerBody{Id: 1},
			expressionsCount: 1,
			tasksCount:       4,
		},
		{
			requestBody:      models.POSTCalculateRequestBody{Expression: "5*(3+10"},
			status:           http.StatusUnprocessableEntity,
			responseBody:     models.POSTCalculateAnswerBody{},
			expressionsCount: 0,
			tasksCount:       0,
		},
		{
			requestBody:      models.POSTCalculateRequestBody{Expression: "7+p"},
			status:           http.StatusUnprocessableEntity,
			responseBody:     models.POSTCalculateAnswerBody{},
			expressionsCount: 0,
			tasksCount:       0,
		},
		{
			requestBody:      models.POSTCalculateRequestBody{Expression: ""},
			status:           http.StatusUnprocessableEntity,
			responseBody:     models.POSTCalculateAnswerBody{},
			expressionsCount: 0,
			tasksCount:       0,
		},
		{
			requestBody:      models.POSTCalculateRequestBody{Expression: "    "},
			status:           http.StatusUnprocessableEntity,
			responseBody:     models.POSTCalculateAnswerBody{},
			expressionsCount: 0,
			tasksCount:       0,
		},
	}

	config := config.NewCfg()
	dbs, _ := NewDBStorage()

	for _, tc := range cases {
		t.Run("Test POSTCalculateHandler", func(t *testing.T) {

			requestBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")
			r := httptest.NewRecorder()

			POSTCalculateHandlerStorage := NewStorage(*config, dbs)
			handler := http.HandlerFunc(POSTCalculateHandler(POSTCalculateHandlerStorage, dbs))
			handler.ServeHTTP(r, req)

			status := r.Code
			if status != tc.status {
				t.Errorf("Error status %v, want %v", status, tc.status)
			}

			if tc.status == http.StatusOK {

				var responseBody models.POSTCalculateAnswerBody
				err = json.NewDecoder(r.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				expressionsCount := len(POSTCalculateHandlerStorage.expressions)
				if expressionsCount != tc.expressionsCount {
					t.Errorf("Error expressions count %v, want %v", expressionsCount, tc.expressionsCount)

				}

				tasksCount := len(POSTCalculateHandlerStorage.tasks)
				if tasksCount != tc.tasksCount {
					t.Errorf("Error tasks count %v, want %v", tasksCount, tc.tasksCount)

				}

			}
		})
	}
}
