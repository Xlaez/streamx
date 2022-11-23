package services

import (
	"context"
	"streamx/models"
	"streamx/responses"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MusicService interface {
	SaveToDb(data models.Music) error
	GetOne(id primitive.ObjectID) (responses.UploadMusic, error)
}

type musicService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewMusicService(col *mongo.Collection, ctx context.Context) MusicService {
	return &musicService{
		col: col,
		ctx: ctx,
	}
}

func (m *musicService) SaveToDb(data models.Music) error {
	_, err := m.col.InsertOne(m.ctx, data)

	if err != nil {
		return err
	}

	return nil
}

func (m *musicService) GetOne(id primitive.ObjectID) (responses.UploadMusic, error) {
	var music responses.UploadMusic

	if err := m.col.FindOne(m.ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&music); err != nil {
		return responses.UploadMusic{}, err
	}

	return music, nil
}
