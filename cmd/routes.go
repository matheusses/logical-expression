package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matheusses/logical-expression/handlers"
)

func setupRoutes(app *fiber.App) {
	app.Get("/expressions", handlers.ListExpressions)
    app.Post("/expressions/:id?", handlers.CreateExpression)
}