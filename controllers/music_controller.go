package controllers

import (
	"errors"
	"net/http"
	config "streamx/configs"
	"streamx/libs"
	"streamx/models"
	"streamx/requests"
	"streamx/services"
	"streamx/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MusicController interface {
	UploadMusic() gin.HandlerFunc
	GetOneMusic() gin.HandlerFunc
	GetMusics() gin.HandlerFunc
	GetMusicsByArtist() gin.HandlerFunc
	DeleteSong() gin.HandlerFunc
}

type musicController struct {
	service services.MusicService
	maker   token.Maker
	config  config.Config
	redis   *redis.Client
}

func NewMusicController(service services.MusicService, maker token.Maker, config config.Config, redis *redis.Client) MusicController {
	return &musicController{
		service: service,
		maker:   maker,
		config:  config,
		redis:   redis,
	}
}

func (c *musicController) UploadMusic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.UploadMusic

		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id := primitive.NewObjectID()

		su, _, e := libs.UploadToCloud(ctx)
		if e != nil {
			ctx.JSON(http.StatusConflict, errorRes(e))
			return
		}

		payload := ctx.MustGet("x-auth-token_payload").(*token.Payload)

		data := models.Music{
			Id:        id,
			Title:     request.Title,
			Cover:     request.Cover,
			User:      payload.Id,
			Artist:    request.Artist,
			CreatedAt: time.Now(),
			File:      su,
		}

		if err := c.service.SaveToDb(data); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusCreated, "uploaded")
	}
}

func (c *musicController) GetOneMusic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.GetOneMusic
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}
		id, err := primitive.ObjectIDFromHex(request.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		music, err := c.service.GetOne(id)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(errors.New("resource not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"music": music})
	}
}

func (c *musicController) GetMusics() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.GetUsers
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		var per_page = (request.Page - 1) * request.Limit

		musics, err := c.service.GetAll(&request.Limit, &per_page)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, musics)
	}
}

func (c *musicController) GetMusicsByArtist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.GetMusicByArtist
		if err := ctx.ShouldBindQuery(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		var per_page = (request.Page - 1) * request.Limit

		musics, err := c.service.GetByArtist(request.Artist, &request.Limit, &per_page)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				ctx.JSON(http.StatusNotFound, errorRes(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, musics)
	}
}

func (c *musicController) DeleteSong() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request requests.GetOneMusic

		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errorRes(err))
			return
		}

		id, err := primitive.ObjectIDFromHex(request.Id)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorRes(err))
			return
		}

		if err = c.service.DeleteSong(id); err != nil {
			ctx.JSON(http.StatusNotFound, errorRes(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "deleted"})
	}
}
