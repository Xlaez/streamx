package services

import (
	"context"
	"errors"
	"streamx/models"
	"streamx/requests"
	"streamx/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(data requests.CreateUser) error
	Login(data requests.LoginUser) (models.User, error)
	VerfiyUser(email string) error
	FindOneByEmail(email string) (models.User, error)
	FindOneById(id primitive.ObjectID) (models.User, error)
	UpdatePassword(email string, password string) error
	UpdateEmail(id primitive.ObjectID, email string) error
	UpdateFields(filter primitive.D, updateObj primitive.D) error
}

type userService struct {
	col *mongo.Collection
	ctx context.Context
}

func NewUserService(col *mongo.Collection, ctx context.Context) UserService {
	return &userService{
		col: col,
		ctx: ctx,
	}
}

func (s *userService) CreateUser(data requests.CreateUser) error {
	id := primitive.NewObjectID()
	hashedPass, err := utils.HashPassword(data.Password)

	if err != nil {
		return err
	}

	filter_by_email := bson.D{{Key: "email", Value: data.Email}}

	var filtered_user models.User

	err = s.col.FindOne(s.ctx, filter_by_email).Decode(&filtered_user)
	if err != mongo.ErrNoDocuments && err != nil {
		return err
	}

	if filtered_user.Email == data.Email {
		er := errors.New("you already have an account, try logging in")
		return er
	}

	user := models.User{
		Id:        id,
		Name:      data.Name,
		Email:     data.Email,
		Password:  hashedPass,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.col.InsertOne(s.ctx, user)

	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(data requests.LoginUser) (models.User, error) {
	filter_user := bson.D{{Key: "email", Value: data.Email}}

	filtered_user := models.User{}

	err := s.col.FindOne(s.ctx, filter_user).Decode(&filtered_user)

	if err != nil {
		return models.User{}, err
	}

	if err = utils.ComparePassword(data.Password, filtered_user.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			new_err := errors.New(" user details is incorrect")
			return models.User{}, new_err
		}
		return models.User{}, err
	}

	return filtered_user, nil
}

func (s *userService) VerfiyUser(email string) error {
	filter := bson.D{{Key: "email", Value: email}}
	updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}}

	_, err := s.col.UpdateOne(s.ctx, filter, updateObj)

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) FindOneByEmail(email string) (models.User, error) {
	filter := bson.D{{Key: "email", Value: email}}

	var user models.User

	if err := s.col.FindOne(s.ctx, filter).Decode(&user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *userService) FindOneById(id primitive.ObjectID) (models.User, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	var user models.User

	if err := s.col.FindOne(s.ctx, filter).Decode(&user); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *userService) UpdatePassword(email string, password string) error {
	hashedPassword, err := utils.HashPassword(password)

	if err != nil {
		return err
	}

	filter := bson.D{{Key: "email", Value: email}}
	updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}}

	_, err = s.col.UpdateOne(s.ctx, filter, updateObj)

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateEmail(id primitive.ObjectID, email string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "email", Value: email}}}}

	_, err := s.col.UpdateOne(s.ctx, filter, updateObj)

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdateFields(filter primitive.D, updateObj primitive.D) error {
	_, err := s.col.UpdateOne(s.ctx, filter, updateObj)

	if err != nil {
		return err
	}

	return nil
}
