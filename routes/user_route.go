package routes

import (
	"streamx/controllers"
	"streamx/token"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, ctl controllers.UserController, tokenMaker token.Maker) {
	// All user routes
	user := router.Group("/api/user").Use(authMiddleWare(tokenMaker))
	authRoute := router.Group("/api/auth")
	authRoute.POST("/register", ctl.CreateUser())
	authRoute.POST("/login", ctl.Login())
	authRoute.POST("/verify", ctl.VerfiyUser())
	authRoute.GET("/reset-password/:email", ctl.GetResetPassword())
	authRoute.PATCH("/reset-password", ctl.ResetPassword())
	user.POST("/change-email", ctl.AskToChangeEmail())
	user.PATCH("/change-email/:digits", ctl.ChangeEmail())
	user.PATCH("/upload/avatar", ctl.UploadAvatar())
	user.GET("/:id", ctl.GetUserById())
	user.DELETE("/", ctl.DeleteAcc())
	user.GET("/many", ctl.GetUsers())
}
