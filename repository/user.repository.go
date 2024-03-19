package repository

import (
	"log"

	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/model/dto"
	"github.com/google/uuid"
)

// Gets User by the username
func GetUserByUsername(username string) (*model.User, int64) {
	var user model.User
	result := db.GetDBConnectionInstance().First(&user, "username = ?", username)

	return &user, result.RowsAffected
}

// Adds new User entity
func AddNewUser(user *model.User) (*model.User, error) {
	result := db.GetDBConnectionInstance().Create(&user)

	if result.Error!=nil {
		return nil, result.Error
	}

	return user, nil
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
	uid, _ := uuid.Parse(id)
	user := model.User{
		Base: model.Base{
			ID: uid,
		},
	}
	result := db.GetDBConnectionInstance().Preload("Images").Find(&user)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return &user, result.RowsAffected
}

func UpdateUser(userDto *dto.UpdateUserDto) (*model.User, int64) {
	uuid, _ := uuid.Parse(userDto.ID)
	user := model.User{
		Base: model.Base{
			ID: uuid,
		},
		Name: userDto.Name,
		Email: userDto.Email,
		Birthday: userDto.Birthday,
	}

	result := db.GetDBConnectionInstance().Model(&user).Updates(user)

	return &user, result.RowsAffected
}

// active value 0 | 1. Field can be updated only via interface{}
func DeactivateUser(user *model.User) (int64, error) {
	result := db.GetDBConnectionInstance().Model(&user).Updates(map[string]interface{}{
			"Active" : 0,
			"SuspendedAt" : user.SuspendedAt,
		})

	return result.RowsAffected, result.Error
}

func ActivateUser(user *model.User) (int64, error) {
	result := db.GetDBConnectionInstance().Model(&user).Updates(map[string]interface{}{
		"Active" : 1,
		"SuspendedAt" : nil,
	})

	return result.RowsAffected, result.Error
}