// Singletone
// MYSQL database connection
package db

import (
	"log"
	"os"
	"sync"

	"github.com/best-nazar/web-app/model"
	"github.com/joho/godotenv"
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
		//Creating single instance now
		lock.Lock()
		defer lock.Unlock()

		if singleInstance == nil {
			err := godotenv.Load(".env")
 			if err != nil{
				panic(".env not found")
			}

			// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
			dsn := os.Getenv("DB_DSN")
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
