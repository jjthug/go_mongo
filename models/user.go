package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name" validate:"required min=2 max=100"`
	Username     *string            `json:"username" validate:"required min=2 max=100"`
	PasswordHash *string            `json:"password_hash" validate:"required"`
	Tags         []*string          `json:"tags"`
	UpdatedAt    time.Time          `json:"updated_at" validate:"required"`
}
