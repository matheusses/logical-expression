package handlers

import (
	"github.com/matheusses/logical-expression/database"
	"github.com/matheusses/logical-expression/models"
	"github.com/gofiber/fiber/v2"
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


func ListExpressions(c *fiber.Ctx) error {
	expressions := []models.LogicalExpression{}
	database.DB.Db.Find(&expressions)

	return c.Status(200).JSON(expressions)
}
