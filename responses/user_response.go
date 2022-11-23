package responses

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserResponse struct {
	Message string `json:"message"`
}

type GetUser struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Avatar    string             `json:"avatar"`
	Verified  bool               `json:"verified"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
