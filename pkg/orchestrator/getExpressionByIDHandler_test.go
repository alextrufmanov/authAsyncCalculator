package orchestrator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

func TestGETExpressionByIDHandler(t *testing.T) {

	cases := []struct {
		expressions []string
		url         string
		status      int
	}{
		{
			expressions: []string{"2+2*2"},
			url:         "/api/v1/expressions",
			status:      http.StatusInternalServerError,
		},
		{
			expressions: []string{"2+2*2"},
			url:         "/api/v1/expressions/1",
			status:      http.StatusInternalServerError,
		},
		{
			expressions: []string{"2+2", "2-2", "2*2", "2/2"},
			url:         "/api/v1/expressions/1",
			status:      http.StatusInternalServerError,
		},
		{
			expressions: []string{},
			url:         "/api/v1/expressions/1",
			status:      http.StatusInternalServerError,
		},
		{
			expressions: []string{"2+2", "2-2", "2*2", "2/2"},
			url:         "/api/v1/expressions/1000",
			status:      http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run("Test GETExpressionByIDHandler", func(t *testing.T) {

			cfg := config.NewCfg()
			dbs, _ := NewDBStorage()
			storage := NewStorage(*cfg, dbs)

			for id, expression := range tc.expressions {
				storage.AppendExpression(int32(id), 1, expression)
			}

			req, err := http.NewRequest("GET", tc.url, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			r := httptest.NewRecorder()
			handler := http.HandlerFunc(GETExpressionByIDHandler(storage, dbs))
			handler.ServeHTTP(r, req)

			status := r.Code
			if status != tc.status {
				t.Errorf("Error status %v, want %v", status, tc.status)
			}

			if status == http.StatusOK {

				var responseBody models.GETexpressionByIDAnswerBody
				err = json.NewDecoder(r.Body).Decode(&responseBody)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}
			}
		})
	}
}
