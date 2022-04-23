// Singletone
// MYSQL database connection
package db

import (
	"log"
	"sync"

	"github.com/best-nazar/web-app/models"
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
			dsn := "testu1:1234@tcp(127.0.0.1:3306)/test_crud?charset=utf8mb4&parseTime=True&loc=Local"
			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			
			if err!=nil {
				log.Fatal(err)
			} else {
				runDbMigration(db)
			}

            singleInstance = db
        } else {
            log.Println("Single instance already created.")
        }
    } else {
        log.Println("Single instance already created.")
    }

    return singleInstance
}

// Migrate the schema to MySQL
func runDbMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.UserActivity{},
	)

	if (err != nil) {
		log.Fatalln(err)
	}
}