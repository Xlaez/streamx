package routes

import (
	"errors"
	"fmt"
	"net/http"
	"streamx/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorixationHeaderKey  = "x-auth-token"
	AuthorizationPayloadKey = "x-auth-token_payload"
)

func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorixationHeader := ctx.GetHeader(AuthorixationHeaderKey)
		if len(authorixationHeader) == 0 {
			err := errors.New("provide an authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		fields := strings.Fields(authorixationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": err.Error()})
			return
		}

		// store payload to context
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
