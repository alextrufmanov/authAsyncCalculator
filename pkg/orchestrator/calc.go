package orchestrator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Ошибки анализа и вычисления арифметического выражения
var (
	ErrEmptyExpression   = errors.New("invalid input, expression is empty")
	ErrDivisionByZero    = errors.New("invalid input, division by zero")
	ErrInvalidExpression = errors.New("invalid input, expression")
)

func ErrUnexpected(item string) error {
	return fmt.Errorf("invalid input, unexpected %s", item)
}

// Поддерживаемые операторы
const operations string = "+-*/"

// Функция определения приоритета указанного оператора
func priority(operation string) int {
	switch operation {
	case "+":
		return 1
	case "-":
		return 1
	case "*":
		return 2
	case "/":
		return 2
	default:
		return -1
	}
}

// Функция подготовки арифметического выражения к преобразованию в RPN,
// разделение на токены
func Split(expression string) ([]string, error) {
	// удаление пробелов
	expression = strings.ReplaceAll(expression, " ", "")
	if len(expression) == 0 {
		// арифметическое выражение оказалось пустым
		return make([]string, 0), ErrEmptyExpression
	}
	// подготовка унарных минусов и плюсов
	if string(expression[0]) == "-" || string(expression[0]) == "+" {
		expression = "0" + expression
	}
	expression = strings.ReplaceAll(expression, "(-", "(0-")
	expression = strings.ReplaceAll(expression, "(+", "(0+")
	// подготовка к разбору на токены
	for _, r := range operations + "()" {
		expression = strings.ReplaceAll(expression, string(r), " "+string(r)+" ")
	}
	// разбор на токены
	return strings.Fields(expression), nil
}

// Функция преобразования арифметического выражения в RPN
func ToRPM(items []string) ([]string, error) {
	var rpm []string
	var stack []string

	for _, item := range items {
		switch {
		case strings.Contains(operations, item):
			for ; len(stack) > 0; stack = stack[:len(stack)-1] {
				if priority(stack[len(stack)-1]) >= priority(item) && stack[len(stack)-1] != "(" {
					rpm = append(rpm, stack[len(stack)-1])
				} else {
					break
				}
			}
			stack = append(stack, item)
		case item == "(":
			stack = append(stack, item)
		case item == ")":
			for ; len(stack) > 0; stack = stack[:len(stack)-1] {
				if stack[len(stack)-1] != "(" {
					rpm = append(rpm, stack[len(stack)-1])
				} else {
					break
				}
			}
			if len(stack) == 0 {
				return make([]string, 0), ErrUnexpected(")")
			}
			stack = stack[:len(stack)-1]
		default:
			rpm = append(rpm, item)
		}
	}

	for ; len(stack) > 0; stack = stack[:len(stack)-1] {
		if stack[len(stack)-1] != "(" {
			rpm = append(rpm, stack[len(stack)-1])
		} else {
			return make([]string, 0), ErrUnexpected("()")
		}
	}

	return rpm, nil
}

// Функция вычисленя арифметического выражения
func Calc(expression string) (float64, error) {
	var stack []float64
	var a, b float64
	// подготавливаем выражение к преобразованию в RPN
	items, err := Split(expression)
	if err != nil {
		return math.NaN(), err
	}
	// преобразуем выражение в RPN
	rpm, err := ToRPM(items)
	if err != nil {
		return math.NaN(), err
	}
	// вычисляем преобразованное в RPN арифметическое выражение
	for _, item := range rpm {
		if strings.Contains(operations, item) {
			// токен RPN оператор, поэтому в стеке
			// уже должны быть минимум 2 числа аргумента
			if len(stack) < 2 {
				return math.NaN(), ErrUnexpected(item)
			}
			// получаем из стека аргументы (два верхних числа)
			a = stack[len(stack)-2]
			b = stack[len(stack)-1]
			stack = stack[:len(stack)-2]
			// выполняем указанную в item операцию и помещаем результат в стек
			switch {
			case item == "+":
				stack = append(stack, a+b)
			case item == "-":
				stack = append(stack, a-b)
			case item == "*":
				stack = append(stack, a*b)
			case item == "/":
				if b == 0 {
					return math.NaN(), ErrDivisionByZero
				}
				stack = append(stack, a/b)
			default:
				return math.NaN(), ErrUnexpected(item)
			}
		} else {
			// токен RPN не оператор (предполагаем, что это число)
			value, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return math.NaN(), ErrUnexpected(item)
			}
			// помещаем это число в стек
			stack = append(stack, value)
		}
	}

	// RPN полностью "разобрана", в стеке должен остаться только
	// результат вычисления арифметического выражения
	if len(stack) != 1 {
		return math.NaN(), ErrInvalidExpression
	}
	return stack[0], nil
}
