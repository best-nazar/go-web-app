// user.activity.repository
package repository

import (
	"log"
	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/models"
)

// Logs new User Activity
func AddUserActivity (activity, data string, userID int64) *models.UserActivity {
	userActivity := models.UserActivity{
		Activity: activity,
		Data: data,
		UserID: uint(userID),
	}

	result := db.GetDBConnectionInstance().Create(&userActivity)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return &userActivity
}