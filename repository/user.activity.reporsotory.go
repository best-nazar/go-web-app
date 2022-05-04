// user.activity.repository
package repository

import (
	"log"
	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
)

// Logs new User Activity
func AddUserActivity (activity, data string, userID int64) *model.UserActivity {
	userActivity := model.UserActivity{
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