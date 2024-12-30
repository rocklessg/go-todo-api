package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID          primitive.ObjectID `bson:"_id"`
	Task        *string            `json:"task" bson:"task" validate:"required,min=3,max=100"`
	Description *string            `json:"description" bson:"description"`
	IsCompleted *bool              `json:"isCompleted" bson:"isCompleted"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
