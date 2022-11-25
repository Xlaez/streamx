package services

import (
	"context"
	"streamx/models"
	"streamx/responses"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MusicService interface {
	SaveToDb(data models.Music) error
	GetOne(id primitive.ObjectID) (responses.UploadMusic, error)
	GetAll(limit *int64, page *int64) ([]responses.UploadMusic, error)
	GetByArtist(artist string, limit *int64, page *int64) ([]responses.UploadMusic, error)
	DeleteSong(id primitive.ObjectID) error
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

func (m *musicService) GetAll(limit *int64, page *int64) ([]responses.UploadMusic, error) {
	cursor, err := m.col.Find(m.ctx, bson.M{}, &options.FindOptions{
		Sort:  bson.D{{Key: "title", Value: 1}},
		Limit: limit,
		Skip:  page,
	})

	if err != nil {
		return nil, err
	}

	var musics []responses.UploadMusic
	cursor.All(m.ctx, &musics)
	return musics, nil
}

// remember to use regex
func (m *musicService) GetByArtist(artist string, limit *int64, page *int64) ([]responses.UploadMusic, error) {
	// searchRegex := primitive.Regex{Pattern: artist, Options: "ig"}.String()
	cursor, err := m.col.Find(m.ctx, bson.D{{Key: "artist", Value: artist}}, &options.FindOptions{
		Sort:  bson.D{{Key: "title", Value: 1}},
		Limit: limit,
		Skip:  page,
	})

	if err != nil {
		return nil, err
	}

	var musics []responses.UploadMusic
	cursor.All(m.ctx, &musics)
	return musics, nil
}

func (m *musicService) DeleteSong(id primitive.ObjectID) error {
	_, err := m.col.DeleteOne(m.ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return err
	}

	return nil
}
