package controllers

import "github.com/gin-gonic/gin"

func errorRes(err error) gin.H {
	return gin.H{"Error": err.Error()}
}
