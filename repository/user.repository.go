package repository

import (
	"log"

	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/models"
)

// Gets User by the username
func GetUserByUsername(username string) (models.User, int64) {
	var user models.User
	result := db.GetDBConnectionInstance().First(&user, "username = ?", username)

	return user, result.RowsAffected
}

// Adds new User entity
func AddNewUser(user *models.User) *models.User {
	result := db.GetDBConnectionInstance().Create(&user)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return user
}