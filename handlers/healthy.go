package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matheusses/logical-expression/database"
)

func Healthy(context *fiber.Ctx) error {

	err := database.DB.Db.Error

	if err != nil {
		return context.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return context.Status(200).JSON(true)
}
