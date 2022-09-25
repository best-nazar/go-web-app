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