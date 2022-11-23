package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Music struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Title     string             `json:"title,omitempty" validate:"required"`
	Cover     string             `json:"avatar" default:"https://cloudinary.com/"`
	Artist    string             `json:"artist" validate:"required"`
	File      string             `json:"file,omitempty" validate:"required"`
	User      primitive.ObjectID `json:"user,omitempty" validate:"required" ref:"users"`
	CreatedAt time.Time          `json:"created_at"`
}
