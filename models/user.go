package models

import (
	"example.com/api-testing/db"
	"github.com/pkg/errors"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save(hashFunc func(string) (string, error)) error {
	query := "INSERT INTO users(email, password) VALUES (?,?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	hashsPassword, err := hashFunc(u.Password)
	result, err := stmt.Exec(u.Email, hashsPassword)

	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials(checkFunc func(string, string) bool) error {
	query := "SELECT id, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrievedpassword string
	err := row.Scan(&u.ID, &retrievedpassword)

	if err != nil {
		return errors.New("Credentials invalid")
	}

	passwordIsValid := checkFunc(u.Password, retrievedpassword)

	if !passwordIsValid {
		return errors.New("Credentials invalid")
	}
	return nil
}
