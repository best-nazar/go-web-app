package repository

import (
	"log"

	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
)

// Gets User by the username
func GetUserByUsername(username string) (*model.User, int64) {
	var user model.User
	result := db.GetDBConnectionInstance().First(&user, "username = ?", username)

	return &user, result.RowsAffected
}

// Adds new User entity
func AddNewUser(user *model.User) *model.User {
	result := db.GetDBConnectionInstance().Create(&user)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return user
}

func GetUsers() []*model.User {
	var users []*model.User
	result := db.GetDBConnectionInstance().Find(&users)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return users
}