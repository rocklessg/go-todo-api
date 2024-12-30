package routes

import (
	controller "go-todo/controllers"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(router *gin.Engine) {
	todoGroup := router.Group("/todos")
	{
		todoGroup.GET("/all-tasks", controller.GetTodos())
		todoGroup.GET("/", controller.GetTodo())
		todoGroup.POST("/add-task", controller.CreateTodo())
		todoGroup.PUT("/", controller.UpdateTodo())
		todoGroup.DELETE("/", controller.DeleteTodo())
	}
}
