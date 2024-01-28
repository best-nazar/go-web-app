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

func FindCasbinRoleById(ID *uint) (*model.CasbinRule, error) {
	var casbinRule model.CasbinRule
	res := db.GetDBConnectionInstance().First(&casbinRule, ID)

	return &casbinRule, res.Error
}

// Adds the binding of Casbin Role to the Username
func AddCasbinUserRole(username string, role string) *model.CasbinRule {
	casbinRule := model.CasbinRule{
		P_type: roleDefinition,
		V0:     username,
		V1:     role,
	}

	db.GetDBConnectionInstance().Create(&casbinRule)

	return &casbinRule
}

func AddCasbinRole(rule *model.CasbinRuleP) model.CasbinRuleP {
	casbinRule := model.CasbinRule{
		P_type: model.GROUP_TYPE_P,
		V0:     rule.V0,
		V1:     rule.V1,
		V2:     rule.V2,
		V3:     rule.V3,
		V4:     rule.V4,
		V5:     rule.V5,
	}

	db.GetDBConnectionInstance().Create(&casbinRule)

	return model.CasbinRuleP(casbinRule)
}

func RemoveCasbinRole(ids []string) error {
	res := db.GetDBConnectionInstance().Delete(model.CasbinRule{}, ids)
	return res.Error
}

func CreateAdminCasbinUserRole(username string) *[]model.CasbinRule {
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

func UpdateCusbinRule(crule *model.CasbinRule) (*model.CasbinRule, error) {
	res := db.GetDBConnectionInstance().Model(&crule).Updates(&crule)

	return crule, res.Error
}
