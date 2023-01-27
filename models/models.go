package models

import (
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type LogicalExpression struct {
	gorm.Model
	Expression string `json:"expression" gorm:"text;not null;default:null`
}

func (self *LogicalExpression) contains(array []string, val string) bool {
	for _, item := range array {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}

func (self *LogicalExpression) performNotOperation(localExpressions []string) []string {
	valueInvert := false
	lsExpresssionNotOperation := []string{}
	for _, value := range localExpressions {
		if value == "!" {
			valueInvert = true
		} else {
			if valueInvert && value == "1" {
				value = "0"
			} else if valueInvert && value == "0" {
				value = "1"
				valueInvert = false
				lsExpresssionNotOperation = append(lsExpresssionNotOperation, value)
			}
		}
		localExpressions = lsExpresssionNotOperation
	}
	return localExpressions
}

func (self *LogicalExpression) performLogicalOperation(localExpressions []string, expressions []string) []string {
	firstArg, _ := strconv.Atoi(localExpressions[0])
	secongArg, _ := strconv.Atoi(localExpressions[2])
	result := 0
	operator := localExpressions[1]
	if operator == "&" {
		result = firstArg & secongArg
	} else {
		result = firstArg | secongArg
	}
	expressions = append(expressions, strconv.Itoa(result))
	return expressions
}

func (self *LogicalExpression) EvaluatePerQueryString(queryString string) (bool, error) {

	expression := self.Expression

	params := strings.Split(queryString, "&")
	for _, value := range params {
		dict := strings.Split(value, "=")
		expression = strings.ToLower(expression)
		expression = strings.ReplaceAll(expression, dict[0], dict[1])
	}

	// Define a map of keyword-symbol associations
	symbols := map[string]string{
		"and": "&",
		"or":  "|",
		"not": "!",
		"(":   "[",
		")":   "]",
		" ":   "",
	}

	// Replace keywords with corresponding symbols
	for keyword, symbol := range symbols {
		expression = strings.ReplaceAll(strings.ToLower(expression), keyword, symbol)
	}

	expressions := []string{}

	// Traversing string from the end
	n := len(expression)
	for i := n - 1; i >= 0; i-- {
		if string(expression[i]) == "[" {
			localExpressions := []string{}
			for len(expressions) > 0 && string(expressions[len(expressions)-1]) != "]" {
				// Creating local expression
				localExpressions = append(localExpressions, expressions[len(expressions)-1])
			}

			// keeping the rest of the expression
			if len(expressions) > 0 {
				expressions = expressions[:len(expressions)-1]
			}

			// Invert the value
			if self.contains(localExpressions, "!") {
				localExpressions = self.performNotOperation(localExpressions)
			}

			// Perform the logical operation
			if len(localExpressions) == 3 {
				expressions = self.performLogicalOperation(localExpressions, expressions)
			}
		} else {
			expressions = append(expressions, string(expression[i]))
		}
	}
	if len(expressions) == 0 {
		err := errors.New("Invalid logical expression")
		return false, err
	}

	result, _ := strconv.Atoi(expressions[len(expressions)-1])
	boolResult := !(result == 0)
	return boolResult, nil
}
