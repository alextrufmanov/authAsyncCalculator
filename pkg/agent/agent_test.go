package agent

import (
	"testing"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
)

func TestCalculate(t *testing.T) {

	cases := []struct {
		task   models.Task
		result bool
		value  float64
	}{
		{
			task: models.Task{
				Id:            1,
				Arg1:          10,
				Arg2:          2,
				Operation:     "+",
				OperationTime: 10,
			},
			result: true,
			value:  12,
		},
		{
			task: models.Task{
				Id:            2,
				Arg1:          10,
				Arg2:          2,
				Operation:     "-",
				OperationTime: 10,
			},
			result: true,
			value:  8,
		},
		{
			task: models.Task{
				Id:            3,
				Arg1:          10,
				Arg2:          2,
				Operation:     "*",
				OperationTime: 10,
			},
			result: true,
			value:  20,
		},
		{
			task: models.Task{
				Id:            4,
				Arg1:          10,
				Arg2:          2,
				Operation:     "/",
				OperationTime: 10,
			},
			result: true,
			value:  5,
		},
		{
			task: models.Task{
				Id:            100,
				Arg1:          10,
				Arg2:          0,
				Operation:     "/",
				OperationTime: 10,
			},
			result: false,
			value:  0,
		},
		{
			task: models.Task{
				Id:            101,
				Arg1:          10,
				Arg2:          2,
				Operation:     "=",
				OperationTime: 10,
			},
			result: false,
			value:  0,
		},
	}

	for _, tc := range cases {
		t.Run("Test GETExpressionByIDHandler", func(t *testing.T) {

			calcResult := calculate(&tc.task)

			if calcResult != tc.result {
				t.Errorf("result %v, want %v", calcResult, tc.result)
			}

			if calcResult && tc.task.Result != tc.value {
				t.Errorf("result value %v, want %v", tc.task.Result, tc.value)
			}
		})
	}
}
