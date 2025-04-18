package dao

import (
	"mini-ups/db"
	"mini-ups/model"
)

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// create user
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}
