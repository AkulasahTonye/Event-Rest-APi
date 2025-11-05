package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "Supersecret"

func GenerateTokens(email string, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userID": userID,
		"Exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	parseToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method")

		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println("Could not parse token")
		return 0, errors.New("Could not parse token.")

	}

	tokenIsValid := parseToken.Valid

	if !tokenIsValid {
		return 0, errors.New("Invalid token!")
	}

	claims, ok := parseToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("Invalid token claims.")
	}

	// email := claims["email"].(string)
	userID, ok := claims["userId"].(int64)

	return userID, nil
}
