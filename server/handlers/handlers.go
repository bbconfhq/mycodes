package handlers

import (
	"github.com/bbconfhq/mycodes/handlers/code"
	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App) {
	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	c := v1.Group("/code")
	c.Get("/", code.GetList)
	c.Post("/", code.Post)
	c.Get("/:uid", code.GetOne)
}
