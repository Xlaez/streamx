package token

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

type Payload struct {
	Id        primitive.ObjectID `json:"id"`
	Email     string             `json:"email"`
	IssuedAt  time.Time          `json:"issued_at"`
	ExpiresAt time.Time          `json:"expired_at"`
}

// Create new payload with fields in {Payload} [struct]
func NewPayload(id primitive.ObjectID, email string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		Id:        id,
		Email:     email,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
