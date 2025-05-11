package orchestrator

import (
	"testing"
)

func TestCalc(t *testing.T) {
	testCases := []struct {
		expression string
		expected   float64
		success    bool
	}{
		{
			expression: "2+2",
			expected:   4,
			success:    true,
		},
		{
			expression: "(2+2)*2",
			expected:   8,
			success:    true,
		},
		{
			expression: "(-2+2)*2",
			expected:   0,
			success:    true,
		},
		{
			expression: "-1/2",
			expected:   -0.5,
			success:    true,
		},
		{
			expression: " 2+ ( -2/( -2 ) + 2)*(-2+2+2) ",
			expected:   8,
			success:    true,
		},
		{
			expression: "-2",
			expected:   -2,
			success:    true,
		},
		{
			expression: "2+",
			expected:   0,
			success:    false,
		},
		{
			expression: "(2*2(2*2)",
			expected:   0,
			success:    false,
		},
		{
			expression: "2-(2+2",
			expected:   0,
			success:    false,
		},
		{
			expression: "2*(2+2)-2)",
			expected:   0,
			success:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run("Test Calc "+testCase.expression, func(t *testing.T) {
			res, err := Calc(testCase.expression)
			if err != nil && testCase.success {
				t.Errorf("Expression: \"%s\". Expected %f, but found error \"%s\".", testCase.expression, testCase.expected, err.Error())
			} else if err == nil && !testCase.success {
				t.Errorf("Expression: \"%s\". Expected error, but found success.", testCase.expression)
			} else if err == nil && testCase.success && res != testCase.expected {
				t.Errorf("Expression: \"%s\". Expected %f, but found %f.", testCase.expression, testCase.expected, res)
			}
		})
	}
}
