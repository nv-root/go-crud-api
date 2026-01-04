package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title" json:"title" validate:"required,min=3"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Category    string             `bson:"category" json:"category" validate:"required"`
	Priority    int                `bson:"priority" json:"priority" validate:"gte=0,lte=5"`
	Status      string             `bson:"status" json:"status" validate:"required,oneof=pending completed cancelled"`
	DueDate     *time.Time         `bson:"due_date" json:"due_date" validate:"omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
