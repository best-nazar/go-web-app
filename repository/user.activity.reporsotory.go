// user.activity.repository
package repository

import (
	"log"

	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
	"github.com/goccy/go-json"
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


// Logs new User Activity
func AddUserActivityData (data interface{}, userID int64) *model.UserActivity {
	strData, derr := json.Marshal(data)
	store := ""
	if derr!=nil {
		store = derr.Error()
	} else {
		store = string(strData)
	}

	userActivity := model.UserActivity{
		Activity: "data",
		Data: store,
		UserID: uint(userID),
	}

	result := db.GetDBConnectionInstance().Create(&userActivity)

	if result.Error!=nil {
		log.Fatal(result.Error)
	}

	return &userActivity
}