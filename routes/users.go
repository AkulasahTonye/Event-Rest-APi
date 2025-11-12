package routes

import (
	"net/http"

	"example.com/api-testing/models"
	"example.com/api-testing/utils"
	"github.com/gin-gonic/gin"
)

// Function variables for mocking in tests
var (
	HashPasswordFunc      = utils.HashPassword
	CheckHashPasswordFunc = utils.CheckHashPassword
	GenerateTokensFunc    = utils.GenerateTokens
	UserSaveFunc          = func(u *models.User) error { return u.Save(HashPasswordFunc) }
	ValidateUserFunc      = func(u *models.User) error { return u.ValidateCredentials(CheckHashPasswordFunc) }
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

	err = UserSaveFunc(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Could not save user",
			"Error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"Message": "User created successfully!"})
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

	err = ValidateUserFunc(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Message": "Could not authenticate user"})
		return
	}

	token, err := GenerateTokensFunc(user.Email, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful!", "token": token})
}
