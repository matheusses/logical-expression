package handlers

import (
	"github.com/matheusses/logical-expression/database"
	"github.com/matheusses/logical-expression/models"
	"github.com/gofiber/fiber/v2"
)

func CreateExpression(c *fiber.Ctx) error {
	expression := new(models.LogicalExpression)
	if err := c.BodyParser(expression); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id := c.SendString(c.Params("id"))
	if  err != nil{
		database.DB.Db.Model(&User{}).Where("id = ?", id).Update("expression", expression.expression)
	}
	else {
		database.DB.Db.Create(&expression)
	}
	return c.Status(200).JSON(expression)
}

func ListExpressions(c *fiber.Ctx) error {
	expressions := []models.LogicalExpression{}
	database.DB.Db.Find(&expressions)

	return c.Status(200).JSON(expressions)
}
