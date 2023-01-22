package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/matheusses/logical-expression/database"
	"github.com/matheusses/logical-expression/models"
)

func CreateExpression(context *fiber.Ctx) error {
	expression := new(models.LogicalExpression)
	if err := context.BodyParser(expression); err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	database.DB.Db.Create(&expression)

	return context.Status(200).JSON(expression)
}

func UpdateExpression(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	expressionModel := &models.LogicalExpression{}

	expression := new(models.LogicalExpression)

	err := context.BodyParser(&expression)

	if err != nil {
		context.Status(fiber.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err

	}

	err = database.DB.Db.Model(expressionModel).Where("id = ?", id).Updates(expression).Error

	if err != nil {
		context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{

			"message": "could not update expression",
		})
		return err

	}

	context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "expression has been successfully updated",
	})

	return nil
}

func DeleteExpression(context *fiber.Ctx) error {
	expressionModel := &models.LogicalExpression{}
	id := context.Params("id")

	if id == "" {
		context.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := database.DB.Db.Delete(expressionModel, id)

	if err.Error != nil {
		context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete expression",
		})
		return err.Error
	}

	context.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "expression has been successfully deleted",
	})

	return nil
}

func ListExpressions(context *fiber.Ctx) error {
	expressions := []models.LogicalExpression{}
	database.DB.Db.Find(&expressions)

	return context.Status(200).JSON(expressions)
}

func EvaluateExpression(context *fiber.Ctx) error {
	querystring := string(context.Request().URI().QueryString())

	id := context.Params("id")
	localExpression := models.LogicalExpression{}
	database.DB.Db.Find(&localExpression, "id = ?", id)

	if localExpression.ID == 0 {
		context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not found expression",
		})
		return nil
	}

	expression := localExpression.Expression

	params := strings.Split(querystring, "&")
	for _, value := range params {
		dict := strings.Split(value, "=")
		expression = strings.ToLower(expression)
		expression = strings.ReplaceAll(expression, dict[0], dict[1])
	}

	result, err := logicalExpressionEvaluation(expression)
	print(err)
	if err != nil {
		context.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": err,
		})
		return nil
	}

	return context.Status(200).JSON(result)
}

func logicalExpressionEvaluation(expression string) (bool, error) {

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
		expression = strings.ReplaceAll(expression, keyword, symbol)
	}

	expressions := []string{}

	// Traversing string from the end
	n := len(expression)
	for i := n - 1; i >= 0; i-- {
		if string(expression[i]) == "[" {
			localExpressions := []string{}
			for len(expressions) > 0 && string(expressions[len(expressions)-1]) != "]" {
				localExpressions = append(localExpressions, expressions[len(expressions)-1])
				expressions = expressions[:len(expressions)-1]
			}

			if len(expressions) > 0 {
				expressions = expressions[:len(expressions)-1]
			}

			// Invert the value
			if contains(localExpressions, "!") {
				valueInvert := false
				lsExpresssionNotOperation := []string{}
				for index := range localExpressions {
					value := localExpressions[index]
					if value == "!" {
						valueInvert = true
					} else {
						if valueInvert && value == "1" {
							value = "0"
						} else if valueInvert && value == "0" {
							value = "1"
						}
						valueInvert = false
						lsExpresssionNotOperation = append(lsExpresssionNotOperation, value)
					}
				}
				localExpressions = lsExpresssionNotOperation
			}

			// Perform the logical operation
			if len(localExpressions) == 3 {
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

func contains(array []string, val string) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}
