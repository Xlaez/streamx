package controllers

import (
	"errors"
	"net/http"
	config "streamx/configs"
	"streamx/libs"
	"streamx/requests"
	"streamx/responses"
	"streamx/services"
	"streamx/token"
	"streamx/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	CreateUser() gin.HandlerFunc
	Login() gin.HandlerFunc
	VerfiyUser() gin.HandlerFunc
	GetResetPassword() gin.HandlerFunc
	ResetPassword() gin.HandlerFunc
	AskToChangeEmail() gin.HandlerFunc
	ChangeEmail() gin.HandlerFunc
	UploadAvatar() gin.HandlerFunc
	GetUserById() gin.HandlerFunc
	GetUsers() gin.HandlerFunc
	DeleteAcc() gin.HandlerFunc
}

type userController struct {
	service services.UserService
	maker   token.Maker
	config  config.Config
	redis   *redis.Client
}

func NewUserController(service services.UserService, maker token.Maker, config config.Config, redis *redis.Client) UserController {
	return &userController{
		service: service,
		maker:   maker,
		config:  config,
		redis:   redis,
	}
}

func (c *userController) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.CreateUser

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err := c.service.CreateUser(requests.CreateUser{
			Name:     request.Name,
			Email:    request.Email,
			Password: request.Password,
		}); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		randInt := utils.RandomIntegers(6)
		err := c.redis.Set(ctx, randInt, request.Email, 0).Err()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"msg": "created", "digits": randInt})
	}
}

func (c *userController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.LoginUser

		if err := ctx.BindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		user, err := c.service.Login(request)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		token, err := createToken(c, user.Id, user.Email, c.config.AccessTokenDuration)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"user": user, "token": token})
	}
}

func (c *userController) VerfiyUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.Verfiy
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		email, err := c.redis.Get(ctx, request.Digits).Result()

		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, errorRes(errors.New("try getting new code")))
			return
		}

		if err = c.service.VerfiyUser(email); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}
		ctx.JSON(http.StatusOK, "verified")
	}
}

func (c *userController) GetResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.AskToReset

		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		randInt := utils.RandomIntegers(6)

		if err := c.redis.Set(ctx, randInt, request.Email, 0).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"digits": randInt})
	}
}

func (c *userController) ResetPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.ResertPassword

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		email, err := c.redis.Get(ctx, request.Digits).Result()

		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, "request for reset digits again")
			return
		}

		user, err := c.service.FindOneByEmail(email)

		if err != nil {
			ctx.JSON(http.StatusNotFound, "cannot find user")
			return
		}

		if err = utils.ComparePassword(request.OldPassword, user.Password); err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				ctx.JSON(http.StatusBadRequest, errorRes(errors.New("password mismatch")))
				return
			}
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		if err = c.service.UpdatePassword(user.Email, request.NewPassword); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"success": true})
	}
}

func (c *userController) AskToChangeEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.AskToReset
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		payload := ctx.MustGet("x-auth-token_payload").(*token.Payload)

		if payload.Email == request.Email {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "please provide a different email, this is your current email"})
			return
		}

		user, err := c.service.FindOneByEmail(request.Email)

		if err != mongo.ErrNoDocuments && user.Email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "email taken, please input another"})
			return
		}

		randInt := utils.RandomIntegers(6)

		if err = c.redis.Set(ctx, randInt, request.Email, 0).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		// send email in the future

		ctx.JSON(http.StatusOK, gin.H{"digits": randInt})
	}
}

func (c *userController) ChangeEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.VerfiyEmail
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		value, err := c.redis.Get(ctx, request.Digits).Result()

		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		payload := ctx.MustGet("x-auth-token_payload").(*token.Payload)

		if err = c.service.UpdateEmail(payload.Id, value); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "updated"})
	}
}

func (c *userController) UploadAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.UploadAvatarReq

		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		su, _, e := libs.UploadToCloud(ctx)

		if e != nil {
			ctx.JSON(http.StatusConflict, errorRes(e))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: id}}
		updateObj := bson.D{{Key: "$set", Value: bson.D{{Key: "avatar", Value: su}}}}

		if err = c.service.UpdateFields(filter, updateObj); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "updated"})
	}
}

func (c *userController) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.GetUserById

		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.ID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		user, err := c.service.FindOneById(id)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		userRes := responses.GetUser{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Verified:  user.Verified,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		ctx.JSON(http.StatusOK, gin.H{"data": userRes})
	}
}

func (c *userController) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.GetUsers
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		var per_page = (request.Page - 1) * request.Limit

		users, err := c.service.GetMany(&request.Limit, &per_page)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("cannot find users")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func (c *userController) DeleteAcc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := ctx.MustGet("x-auth-token_payload").(*token.Payload)

		user, err := c.service.FindOneById(payload.Id)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if user.Avatar != "" {
			if err = libs.DeleteFromCloud(user.Avatar, ctx); err != nil {
				ctx.JSON(http.StatusExpectationFailed, errorRes(err))
				return
			}
		}

		if err = c.service.DeleteAcc(user.Id); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, "successful")
	}
}

func createToken(s *userController, id primitive.ObjectID, email string, duration time.Duration) (string, error) {
	accToken, err := s.maker.CreateToken(id, email, duration)

	if err != nil {
		return "", err
	}
	return accToken, nil
}
