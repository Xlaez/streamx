package token

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Maker interface {
	CreateToken(id primitive.ObjectID, email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
