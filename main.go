package main

import (
	"os"

	"github.com/gin-gonic/gin"
	routes "go-todo/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.TodoRoutes(router)

	router.Run(":" + port)
}
