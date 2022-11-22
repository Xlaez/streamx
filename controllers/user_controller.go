package controllers

import (
	"errors"
	"net/http"
	config "streamx/configs"
	"streamx/requests"
	"streamx/services"
	"streamx/token"
	"streamx/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	CreateUser() gin.HandlerFunc
	Login() gin.HandlerFunc
	VerfiyUser() gin.HandlerFunc
	GetResetPassword() gin.HandlerFunc
	ResetPassword() gin.HandlerFunc
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

func createToken(s *userController, id primitive.ObjectID, email string, duration time.Duration) (string, error) {
	accToken, err := s.maker.CreateToken(id, email, duration)

	if err != nil {
		return "", err
	}
	return accToken, nil
}
