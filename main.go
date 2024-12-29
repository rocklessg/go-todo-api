package main

import(
	"os"

	routes "go-todo/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.TodoRoutes(router)

	router.Run(":" + port)
}