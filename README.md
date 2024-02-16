# Go gin app

Inspired by [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin).

# WEB-SERVICE UI/API
- You can register user and have a back-office Admin page based on the user role.
- Service supports three roles (or groups): *guest*, *memeber* and *admin*
- You can configure groups and access level to URL resources

1. An authorization middleware: https://github.com/gin-contrib/authz
2. Auth policy editor: https://casbin.org/editor/
3. Auth configuration: authz_model.conf & authz_policy.csv
4. In Casbin, an access control model is abstracted into a CONF file based on the PERM metamodel (Policy, Effect, Request, Matchers) https://github.com/casbin/casbin

5. Default the policy rules: see ``automigration()`` in ``dbconnection.go``
- Restart the App to reload rules.
6. Groups hierarchy
- admin
- - member
- - - guest

6. Admin user
- To register a User with admin rights, change the configuration in ``config.yaml`` to the next 
```
default_casbin_group: admin
```
- Then hit the ``/u/register`` URL and fill in the form. Revert back the configuration in ``config.yaml`` to 'member'.

7. API:
- Authorization - Basic
- Username - guest
- password - leave empty
- Headers:
	``Accept:application/json``
or
	``Accept:application/xml``

Note: Make sure there is no double slashed in URL like: //u/login (must be /u/login). 403 error will be caused by Casbin.

# Code conventions:

1. Response:
```
	Render(c, gin.H{
		"title": "Acount is locked ",
		"user":  user,
	}, "register-successful.html", http.StatusOK)
```
2. Error handling:
- Add the binding to the Model:
```
type CasbinRole struct {
	ID     			uint   		`gorm:"primaryKey"`
	Title  			string 		`json:"title" gorm:"index" form:"title" binding:"alphanum,min=3"`
	InheritedFrom 	string		`json:"inheritedFrom" gorm:"index" form:"inheritedFrom" binding:"excludesall= "`
}
```
- Store error:
```
	c.Error(errors.New("Not found"))
```
- Return error:
```
if len(c.Errors) > 0 {
	Render(c, gin.H{
		"title": "Error",
		"description": "Account is locked",
		"errors" : helpers.Errors(c),
	}, "errors.html", http.StatusBadRequest)
}
```

8. User Activity Data collecting
To store the history of user actions like: page opened, data added/updated/deleted
Turn on/off this option in config.yaml and check DB table 'user_activities'
```
user_activity_logging: true
```

9. DB Migration
- Automigration: see "dbconnection.go"
- Condition - if no tables created then auto run migration. In other words, drop all tables in schema, run the server, and hit any endpoint.

10. APP configuration extention
- Add the propoerty to config.yaml
- Add the property to the Struct 'model.Config'
- Read the config
```
	config := c.MustGet("config").(model.Config)
```