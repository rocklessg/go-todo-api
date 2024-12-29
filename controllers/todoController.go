package controllers

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"go-todo/database"
	"go-todo/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection *mongo.Collection = database.OpenCollection(database.TodoDbContext(), "todo")
var validateStruct = validator.New()

// CreateTodo creates a new todo
func CreateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var todo models.Todo

		// Bind the JSON to the todo model
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the todo model
		if err := validateStruct.Struct(todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo.ID = primitive.NewObjectID()
		isCompleted := false
		todo.IsCompleted = &isCompleted
		todo.CreatedAt = time.Now()
		todo.UpdatedAt = time.Now()

		// Insert the todo record into the db
		result, err := todoCollection.InsertOne(ctx, todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating Todo"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

// GetTodos returns all todos with pagination
func GetTodos() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get page and limit query parameters
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit < 1 {
			limit = 10
		}

		// Calculate skip value
		skip := (page - 1) * limit

		// Find options with limit and skip
		findOptions := options.Find()
		findOptions.SetLimit(int64(limit))
		findOptions.SetSkip(int64(skip))

		// Find todos
		cursor, err := todoCollection.Find(ctx, bson.M{}, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching Todos"})
			return
		}
		defer cursor.Close(ctx)

		var todos []models.Todo
		if err = cursor.All(ctx, &todos); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while decoding (processing) Todo records"})
			return
		}

		// Get total count of todos
		total, err := todoCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while counting Todos"})
			return
		}

		// Calculate total pages
		totalPages := int(math.Ceil(float64(total) / float64(limit)))

		c.JSON(http.StatusOK, gin.H{
			"data":       todos,
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		})
	}
}

// GetTodo returns a single todo by ID
func GetTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var todo models.Todo

		// Get the ID from the URL (query param)and check if exist in the db
		id, err := primitive.ObjectIDFromHex(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// fetch the todo with matching ID
		err = todoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching Todo"})
			return
		}
		c.JSON(http.StatusOK, todo)
	}
}

// UpdateTodo updates a todo by ID
func UpdateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var todo models.Todo

		// Get the ID from the URL (query param) and check if it is valid
		id, err := primitive.ObjectIDFromHex(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Bind the JSON to the todo model
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Build the update object dynamically
		update := bson.M{"$set": bson.M{}}
		if todo.Task != nil {
			update["$set"].(bson.M)["task"] = *todo.Task
		}
		if todo.Description != nil {
			update["$set"].(bson.M)["description"] = *todo.Description
		}
		if todo.IsCompleted != nil {
			update["$set"].(bson.M)["isCompleted"] = *todo.IsCompleted
		}

		// Always update the "updatedAt" field
		update["$set"].(bson.M)["updatedAt"] = time.Now()

		// If no fields are updated, return a bad request
		if len(update["$set"].(bson.M)) == 1 { // Only "updatedAt" was added
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		// Update the todo record in the db
		result, err := todoCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating Todo"})
			return
		}

		// Check if the record was found and updated
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo updated successfully"})
	}
}

// DeleteTodo deletes a todo by ID
func DeleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Get the ID from the URL (query param) and check if it is valid
		id, err := primitive.ObjectIDFromHex(c.Query("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Delete the todo record from the db
		result, err := todoCollection.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting Todo"})
			return
		}

		// Check if the record was found and deleted
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}
