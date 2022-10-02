# Go gin app

This is the code from the article [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin).


1. An authorization middleware: https://github.com/gin-contrib/authz
2. Auth policy editor: https://casbin.org/editor/
3. Auth configuration: authz_model.conf & authz_policy.csv
4. In Casbin, an access control model is abstracted into a CONF file based on the PERM metamodel (Policy, Effect, Request, Matchers) https://github.com/casbin/casbin

5. Adding the policy rule:
- Run query: ```INSERT INTO casbin_rule (p_type, v0, v1, v2) VALUES ('p, guest, /article/view/*, GET);```
- Restart the App to reload rules.
6. Adding the group:
- Run query: ```INSERT INTO casbin_rule (p_type, v0, v1) VALUES ('g', 'test1', 'guest');```
- Restart the App to reload rules. More examples at ``__authz_policy.csv``

7. API:
Authorization - Basic
Username - guest
password - leave empty
Headers:
	Accept:application/json
or
	Accept:application/xml

Note: Make sure there is no double slashed in URL like: //u/login (must be /u/login). 403 error will be caused by Casbin.

##Code conventions:

1. Response:
```
	Render(c, gin.H{
		"title": "Users and Roles",
		"payload": casbins,
		"admin-dashboard.html",
		http.StatusOk)
```
2. Error handling:
2.1. Add the binding to the Model:
```
type CasbinRole struct {
	ID     			uint   		`gorm:"primaryKey"`
	Title  			string 		`json:"title" gorm:"index" form:"title" binding:"alphanum,min=3"`
	InheritedFrom 	string		`json:"inheritedFrom" gorm:"index" form:"inheritedFrom" binding:"excludesall= "`
}
```
2.2. Check the errorMsgHandl.go, Make sure the binding tag is in switch/case.
2.3. In controlller:
2.3.1 Validate the request:
```
if err != nil {
	casbins := repository.GetGroupRoles()
	errView := errorSrc.MakeErrorView("Add role", err)

	err := []string{err.Error()}
	errView := errorSrc.MakeErrorView("Add role", err)
		Render(c, gin.H{
			"payload": data,
			"error": errView},
			"template.html",
			http.StatusBadRequest)
}			
```
Note: 
3. Groups hierarchy
- guest
- - member
- - - admin