package middlewares

import (
	"net/http"

	"example.com/api-testing/utils"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "userId"

var verifyTokenFunc = utils.VerifyToken

func Authenticate(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized."})
		return
	}

	userId, err := verifyTokenFunc(token)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized."})
		return
	}
	ctx.Set(UserIDKey, userId)
	ctx.Next()
}
