package routes

import (
	"streamx/controllers"
	"streamx/token"

	"github.com/gin-gonic/gin"
)

func MusicRoutes(router *gin.Engine, ctl controllers.MusicController, tokenMaker token.Maker) {
	// All user routes
	upload := router.Group("/api/music").Use(authMiddleWare(tokenMaker))
	upload.POST("/new", ctl.UploadMusic())
	upload.GET("/single/:id", ctl.GetOneMusic())
}
