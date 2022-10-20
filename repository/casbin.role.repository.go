// casbin.role.repository
package repository

import (
	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
)

func GetGroupRoles() []model.CasbinRole {
	var groupRoles = []model.CasbinRole{}

	db.GetDBConnectionInstance().Find(&groupRoles)

	return groupRoles
}

func SaveCasbinRole(casbinRole *model.CasbinRole) model.CasbinRole {
	db.GetDBConnectionInstance().Create(&casbinRole)

	return *casbinRole
}

func DeleteCasbinRole(casbinRole *model.CasbinRole) (int64, error) {
	res := db.GetDBConnectionInstance().Delete(&casbinRole)
	
	return res.RowsAffected, res.Error
}