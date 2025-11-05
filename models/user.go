package models

import (
	"example.com/api-testing/db"
	"example.com/api-testing/utils"
	"github.com/pkg/errors"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashsPassword, err := utils.HashPassword(u.Password)
	result, err := stmt.Exec(u.Email, hashsPassword)

	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedpassword string
	err := row.Scan(&u.ID, &retrievedpassword)

	if err != nil {
		return errors.New("Credentials invalid")
	}

	passwordIsValid := utils.CheckHashPassword(u.Password, retrievedpassword)

	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}
	return nil
}
