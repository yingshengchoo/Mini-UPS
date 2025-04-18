package service

import (
	"fmt"
	"log"
	"mini-ups/dao"
	"mini-ups/model"

	"golang.org/x/crypto/bcrypt"
)

func GetUserByUsername(username string) (*model.User, error) {
	return dao.GetUserByUsername(username)
}

func RegisterUser(username, password string) error {
	// check if user exists
	user, _ := dao.GetUserByUsername(username)
	if user != nil {
		log.Printf("Username <%s> already exists!", username)
		return fmt.Errorf("username <%s> already exists", username)
	}

	// encrypt password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create new user
	newUser := model.User{
		Username: username,
		Password: string(hashedPwd),
	}

	return dao.CreateUser(&newUser)
}
