package models

import (
	"errors"
	"regexp"
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

	if len(expressions) > 0 {
		expressions = append(expressions[:len(expressions)-1], strconv.Itoa(result))
	} else {
		expressions = append(expressions, strconv.Itoa(result))
	}
	return expressions
}

func (self *LogicalExpression) convertQueryStringToExpression(queryString string) (string, error) {
	expression := self.Expression
	params := strings.Split(queryString, "&")
	// Base to interepret the string in (default is 10)
	base := 10
	// Size (in bits) of the resulting integer type
	bitSize := 32
	pattern := ".*[a-zA-Z].*"

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

	// Replay query string to value
	for _, value := range params {
		dict := strings.Split(value, "=")
		value, err := strconv.ParseInt(dict[1], base, bitSize)
		if err != nil || (value != 0 && value != 1) {
			err := errors.New("Invalid value in query string: " + dict[0] + "=" + dict[1])
			return expression, err
		}
		expression = strings.ToLower(expression)
		expression = strings.ReplaceAll(expression, dict[0], dict[1])
	}

	match, _ := regexp.MatchString(pattern, expression)
	if match {
		err := errors.New("Paramaters not fill")
		return expression, err
	}

	if len(expression) > 0 && expression[0] != '[' {
		expression = "[" + expression + "]"
	}

	return expression, nil
}

func (self *LogicalExpression) EvaluatePerQueryString(queryString string) (bool, error) {

	expression, err := self.convertQueryStringToExpression(queryString)

	if err != nil {
		return false, err
	}

	expressions := []string{}

	// Traversing string from the end
	n := len(expression)
	for i := n - 1; i >= 0; i-- {
		if string(expression[i]) == "[" {
			localExpressions := []string{}
			// Solving expression - While the logical expression is solving and does not reach the final based on a close bracket
			for len(expressions) > 0 && string(expressions[len(expressions)-1]) != "]" {
				// Creating local expression
				localExpressions = append(localExpressions, expressions[len(expressions)-1])
				// Updating expressions decrementing local expression
				expressions = expressions[:len(expressions)-1]

				// Invert the value
				if self.contains(localExpressions, "!") {
					localExpressions = self.performNotOperation(localExpressions)
				}

				// Perform the logical operation
				if len(localExpressions) == 3 {
					expressions = self.performLogicalOperation(localExpressions, expressions)
					localExpressions = []string{}
				}
			}
			if len(localExpressions) > 0 && len(expressions) == 0 {
				expressions = localExpressions
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
