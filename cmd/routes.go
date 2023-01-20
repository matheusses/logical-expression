package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matheusses/logical-expression/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/expressions", handlers.ListExpressions)
    app.Post("/expressions", handlers.CreateExpression)
	app.Put("/expressions/:id", handlers.UpdateExpression)
	app.Delete("/expressions/:id", handlers.DeleteExpression)
}