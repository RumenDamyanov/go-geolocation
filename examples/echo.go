//go:build ignore
// +build ignore

// This file is for documentation/example purposes and is not built with the main module.

package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rumendamyanov/go-geolocation/echoadapter"
)

func main() {
	e := echo.New()
	e.Use(echoadapter.Middleware())
	e.GET("/", func(c echo.Context) error {
		loc := echoadapter.FromContext(c)
		return c.JSON(200, loc)
	})
	e.Logger.Fatal(e.Start(":8082"))
}
