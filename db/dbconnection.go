// Singletone
// MYSQL database connection
package db

import (
	"log"
	"sync"

	"github.com/best-nazar/web-app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// The singleton struct must return the same instance
// whenever multiple goroutines are trying to access that instance.
var lock = &sync.Mutex{}

var singleInstance *gorm.DB

// Gets DB connection instance
func GetDBConnectionInstance() *gorm.DB {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			log.Println("Creating single instance now.")

			// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
			dsn := "testu1:1234@tcp(localhost:3306)/test_crud?charset=utf8mb4&parseTime=True&loc=Local"
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

			if err == nil {
				runDbMigration(db)
				insertInitData(db)
			} else {
				log.Fatal(err)
			}

			singleInstance = db
		} else {
			log.Println("Single instance already created 1.")
		}
	} else {
		log.Println("Single instance already created 2.")
	}

	return singleInstance
}

// Migrate the schema to MySQL
func runDbMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.UserActivity{},
		&model.CasbinRule{},
		&model.CasbinRole{},
	)

	if err != nil {
		log.Fatalln(err)
	}
}

// Insert initial data to the DB
func insertInitData(db *gorm.DB) {
	casbinRole := model.CasbinRole{}
	res := db.First(&casbinRole)

	if res.Error != nil {
		db.Model(&model.CasbinRole{}).Create([]map[string]interface{}{
			{"title": model.GUEST_ROLE, "IsSystem": true, "Description": "limited access"},
			{"title": model.USER_ROLE, "IsSystem": true, "Description": "common role"},
			{"title": model.ADMIN_ROLE, "IsSystem": true, "Description": "system administrator"},
		})
		// Setting hireracy (Admin is in User, User is in Guest)
		// Setting default routes level of access
		db.Model(&model.CasbinRule{}).Create([]map[string]interface{}{
			{"P_type": model.GROUP_TYPE_G, "V0": model.ADMIN_ROLE, "V1": model.USER_ROLE, "V2": "", "V3": "", "V4": "", "V5": ""},
			{"P_type": model.GROUP_TYPE_G, "V0": model.USER_ROLE, "V1": model.GUEST_ROLE, "V2": "", "V3": "", "V4": "", "V5": ""},
			{"P_type": model.GROUP_TYPE_P, "V0": model.GUEST_ROLE, "V1": "/", "V2": "GET", "V3": "", "V4": "", "V5": ""},
			{"P_type": model.GROUP_TYPE_P, "V0": model.GUEST_ROLE, "V1": "/u/*", "V2": "*", "V3": "", "V4": "", "V5": ""},
			{"P_type": model.GROUP_TYPE_P, "V0": model.ADMIN_ROLE, "V1": "/admin/*", "V2": "*", "V3": "", "V4": "", "V5": ""},
		})
	}
}
