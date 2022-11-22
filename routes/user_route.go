package routes

import (
	"streamx/controllers"
	"streamx/token"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, ctl controllers.UserController, tokenMaker token.Maker) {
	// All user routes
	auth := router.Group("/api/auth").Use(authMiddleWare(tokenMaker))
	authRoute := router.Group("/api/auth")
	authRoute.POST("/register", ctl.CreateUser())
	authRoute.POST("/login", ctl.Login())
	authRoute.POST("/verify", ctl.VerfiyUser())
	authRoute.GET("/reset-password/:email", ctl.GetResetPassword())
	authRoute.PATCH("/reset-password", ctl.ResetPassword())
	auth.POST("/change-email", ctl.AskToChangeEmail())
	auth.PATCH("/change-email/:digits", ctl.ChangeEmail())
	auth.PATCH("/upload/avatar", ctl.UploadAvatar())

}
