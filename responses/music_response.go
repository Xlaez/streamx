package responses

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadMusic struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     string             `json:"title"`
	Cover     string             `json:"cover"`
	File      string             `json:"file"`
	User      primitive.ObjectID `json:"user"`
	Artist    string             `json:"artist" required:"true"`
	CreatedAt time.Time          `json:"created_at"`
}
