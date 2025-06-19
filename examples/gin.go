//go:build ignore
// +build ignore

// This file is for documentation/example purposes and is not built with the main module.

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rumendamyanov/go-geolocation/ginadapter"
)

func main() {
	r := gin.Default()
	r.Use(ginadapter.Middleware())
	r.GET("/", func(c *gin.Context) {
		loc := ginadapter.FromContext(c)
		c.JSON(200, loc)
	})
	r.Run(":8081")
}
