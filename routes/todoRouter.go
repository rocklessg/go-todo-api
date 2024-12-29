package routes

import (
	"go-todo/controllers"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(router *gin.Engine) {
	todoGroup := router.Group("/todos")
	{
		todoGroup.GET("/", controllers.GetTodos())
		todoGroup.GET("/:id", controllers.GetTodo())
		todoGroup.POST("/", controllers.CreateTodo())
		todoGroup.PUT("/:id", controllers.UpdateTodo())
		todoGroup.DELETE("/:id", controllers.DeleteTodo())
	}
}