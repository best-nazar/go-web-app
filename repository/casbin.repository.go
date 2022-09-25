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