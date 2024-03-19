package repository

import (
	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
)

func SaveImage(img *model.Image) error {
	tx := db.GetDBConnectionInstance().Save(&img)
	
	return tx.Error
}

func DeleteImage(img *model.Image) error {
	tx := db.GetDBConnectionInstance().Delete(&img)
	
	return tx.Error
}
