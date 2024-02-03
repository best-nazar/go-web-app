package repository

import (
	"log"

	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
	"gorm.io/gorm"
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

func FindUserById(id string) (*model.User, int64) {
	var user *model.User
	result := db.GetDBConnectionInstance().Find(&user, id)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return user, result.RowsAffected
}

func UpdateUser(user *model.UpdateUser) (*model.User, int64) {
	var result *gorm.DB

	lookupUser := model.User{
		ID: user.ID,
	}
	
	findResult := db.GetDBConnectionInstance().First(&lookupUser)
	
	if findResult.Error!=nil {
		log.Fatal(findResult.Error)
	}

	if findResult.RowsAffected > 0 {
		result = db.GetDBConnectionInstance().Model(&lookupUser).Updates(model.User{
			Name: user.Name,
			Email: user.Email,
			Birthday: user.Birthday,
			Active: user.Active,
		})
	}

	return &lookupUser, result.RowsAffected
}