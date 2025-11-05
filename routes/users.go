package routes

import (
	"net/http"

	"example.com/api-testing/models"
	"example.com/api-testing/utils"
	"github.com/gin-gonic/gin"
)

func signup(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message": "Could not parse request data",
			"Error":   err.Error(),
		})
		return
	}
	err = user.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Could not save user",
			"Error":   err.Error(),
		})
		ctx.JSON(http.StatusCreated, gin.H{"Message": "User created successfully!"})
	}
}

func login(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message": "Could not parse request data",
			"Error":   err.Error(),
		})
		return
	}
	err = user.ValidateCredentials()

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Message": "Could not authenticate user"})
		return
	}

	token, err := utils.GenerateTokens(user.Email, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not authenticate user."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
