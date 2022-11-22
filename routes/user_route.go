package routes

import (
	"streamx/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, ctl controllers.UserController) {
	// All user routes
	auth := router.Group("/api/auth")
	auth.POST("/register", ctl.CreateUser())
	auth.POST("/login", ctl.Login())
	auth.POST("/verify", ctl.VerfiyUser())
	auth.GET("/reset-password/:email", ctl.GetResetPassword())
	auth.PATCH("/reset-password", ctl.ResetPassword())

}
