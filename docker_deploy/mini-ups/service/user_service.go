package service

import (
	"mini-ups/dao"
	"mini-ups/model"
)

func GetUserByUsername(username string) (*model.User, error) {
	return dao.GetUserByUsername(username)
}
