package service

import (
	"fmt"
	"log"
	"mini-ups/dao"
	"mini-ups/model"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func LoginUser(username, password string, c *gin.Context) error {
	// check if user exists
	user, _ := dao.GetUserByUsername(username)
	if user == nil {
		log.Printf("Username <%s> not exists!", username)
		return fmt.Errorf("username <%s> not exists", username)
	}

	// match the info
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// log.Printf("err: %v", err)

	if err != nil {
		// not match
		return fmt.Errorf("wrong username or password")
	}

	// matched
	log.Println("match successfully")
	session := sessions.Default(c)
	session.Set("user", username)
	err = session.Save()
	if err != nil {
		log.Println("Failed to save session:", err)
	}

	return nil

}
