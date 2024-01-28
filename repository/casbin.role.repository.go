// casbin.role.repository
package repository

import (
	"errors"

	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
)

// Get full list of the Roles for Casbin RBAC
func ListRoles() *[]model.CasbinRole {
	var groupRoles = []model.CasbinRole{}
	db.GetDBConnectionInstance().Find(&groupRoles)

	return &groupRoles
}

func SaveCasbinRole(casbinRole *model.CasbinRole) model.CasbinRole {
	db.GetDBConnectionInstance().Create(&casbinRole)

	return *casbinRole
}

func DeleteCasbinRole(casbinRole *model.CasbinRole) (int64, error) {
	rs := db.GetDBConnectionInstance().First(&casbinRole)

	if casbinRole.IsSystem {
		return 0, errors.New("the Role is System")
	}

	if rs.RowsAffected > 0 {
		var group model.CasbinRule
		result := db.GetDBConnectionInstance().Where("v1 = ?", casbinRole.Title).Find(&group)

		if result.RowsAffected > 0 {
			return 0, errors.New("the Role is in use")
		}
	}
	
	res := db.GetDBConnectionInstance().Delete(&casbinRole)
	
	return res.RowsAffected, res.Error
}

// Deletes Users assigned to the Casbin group.
func DeleteCasbinGroup (list *model.UsersList) (int64, error) {
	casbinRules := []model.CasbinRule{}
	res := db.GetDBConnectionInstance().Where("V0 IN ? AND p_type = ? and V1 = ?", list.Users, model.GROUP_TYPE_G, list.Group).Delete(&casbinRules)

	return res.RowsAffected, res.Error
}

func FindCasbinRolebyName(title string) (model.CasbinRole, error) {
	var role model.CasbinRole
	res := db.GetDBConnectionInstance().First(&role, "title=?", title)

	return role, res.Error
}