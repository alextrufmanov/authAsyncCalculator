package orchestrator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alextrufmanov/asyncCalculator/pkg/config"
	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

func TestGETExpressionsHandler(t *testing.T) {

	cases := []struct {
		expressions []string
		status      int
	}{
		{
			expressions: []string{"2+2*2"},
			status:      http.StatusOK,
		},
		{
			expressions: []string{"2+2", "2-2", "2*2", "2/2"},
			status:      http.StatusOK,
		},
		{
			expressions: []string{},
			status:      http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run("Test GETExpressionsHandler", func(t *testing.T) {

			cfg := config.NewCfg()
			dbs, _ := NewDBStorage()
			storage := NewStorage(*cfg, dbs)

			for id, expression := range tc.expressions {
				storage.AppendExpression(int32(id), 0, expression)
			}

			req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			r := httptest.NewRecorder()
			handler := http.HandlerFunc(GETExpressionsHandler(storage, dbs))
			handler.ServeHTTP(r, req)

			status := r.Code
			if status != tc.status {
				t.Errorf("Error status %v, want %v", status, tc.status)
			}

			if status == http.StatusOK {

				var responseBody models.GETexpressionsAnswerBody
				err = json.NewDecoder(r.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				responseCount := len(responseBody.Expressions)
				wantCount := len(tc.expressions)
				if responseCount != wantCount {
					t.Errorf("responseBody.Expressions count %v, want %v.", responseCount, wantCount)
				}

			}
		})
	}
}
