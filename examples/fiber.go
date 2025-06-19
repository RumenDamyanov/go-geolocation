//go:build ignore
// +build ignore

// This file is for documentation/example purposes and is not built with the main module.

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rumendamyanov/go-geolocation/fiberadapter"
)

func main() {
	app := fiber.New()
	app.Use(fiberadapter.Middleware())
	app.Get("/", func(c *fiber.Ctx) error {
		loc := fiberadapter.FromContext(c)
		return c.JSON(200, loc)
	})
	app.Listen(":8083")
}
