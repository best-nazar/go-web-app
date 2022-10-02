// casbin.repository
package repository

import (
	"github.com/best-nazar/web-app/db"
	"github.com/best-nazar/web-app/model"
)
// table name for RBAC
const (
	//Policy definition
	//p = sub, obj, act
	policyDefinition = "p"
	//role_definition
	//g = _, _
	roleDefinition = "g"
)

func GetCasbinPolicies() []model.CasbinRule {
	var casbinPolicies []model.CasbinRule

	db.GetDBConnectionInstance().Find(&casbinPolicies, "p_type=?", policyDefinition)

	return casbinPolicies
}

func GetCasbinRoles() []model.CasbinRule {
	var casbinRoles []model.CasbinRule

	db.GetDBConnectionInstance().Find(&casbinRoles, "p_type=?", roleDefinition)

	return casbinRoles
}

// Adds the binding of Casbin Role to the Username
func AddCasbinUserRole(username string, role string) *model.CasbinRule {
	casbinRule := model.CasbinRule{
		P_type: roleDefinition,
		V0: username,
		V1: role,
	}

	db.GetDBConnectionInstance().Create(&casbinRule)

	return &casbinRule
}

func SaveAdminCasbinUserRole(username string) *[]model.CasbinRule {
	var casbinRules = []model.CasbinRule{
		{P_type: roleDefinition,
		V0: username,
		V1: model.ADMIN_ROLE},
		{P_type: roleDefinition,
		V0: username,
		V1: model.USER_ROLE},
		{P_type: roleDefinition,
		V0: username,
		V1: model.GUEST_ROLE},
	}

	db.GetDBConnectionInstance().Create(&casbinRules)

	return &casbinRules
}