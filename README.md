# Go Gin app

Inspired by [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin).

# WEB-SERVICE UI/API
- You can register a User and have a back-office (Admin page) based on the User role.
- Service supports three roles (or groups): *guest*, *member* and *admin*
- You can configure groups and access levels to URL resources

## Resources
1. An authorization middleware: https://github.com/gin-contrib/authz
2. Auth policy editor: https://casbin.org/editor/
3. Auth configuration: authz_model.conf & authz_policy.csv
4. In Casbin, an access control model is abstracted into a CONF file based on the PERM metamodel (Policy, Effect, Request, Matchers) https://github.com/casbin/casbin

## Initial configuration
1. Storage of Default policy rules: see ``automigration()`` in ``dbconnection.go``
- Restart the App to reload rules.
2. Group hierarchy
- admin
- - member
- - - guest
3. Admin user
- To register a User with admin rights, change the configuration in ``config.yaml`` to the next 
```
default_casbin_group: admin
```
- Then hit the ``/u/register`` URL and fill in the form. Revert back the configuration in ``config.yaml`` to 'member'.
- Unlock the Admin User by setting columnn 'active' = 1 in users table in MySQL DB.

7. API:
- Authorization - Basic
- Username - guest
- password - leave empty
- Headers:
	``Accept:application/json``
or
	``Accept:application/xml``

Note: Make sure there is no double slashed in URL like: //u/login (must be /u/login). 403 error will be caused by Casbin.

# Code convention examples:

1. Response:
```
	Render(c, gin.H{
		"title": "Acount is locked ",
		"data":  data,
	}, "register-successful.html", http.StatusOK)
```
2. Error handling:
- Add the binding to the Model:
```
type CasbinRole struct {
	Title  			string 		`json:"title" gorm:"index" form:"title" binding:"alphanum,min=3"`
	InheritedFrom 	string		`json:"inheritedFrom" gorm:"index" form:"inheritedFrom" binding:"excludesall= "`
}
```
- Store error in context:
```
	c.Error(fmt.Errorf("Not found"))
```


3. User Activity Data collecting
To store the history of user actions like: page opened, data added/updated/deleted
Turn on/off this option in config.yaml and check DB table 'user_activities'
```
user_activity_logging: true
```

4. DB Migration
- Auto migration: add more structs to "dbconnection.go"
- Condition - if DB schema is empty the APP will run the migration on the first URL hit (any URL). In other words, drop all tables in schema, run the server, and hit any endpoint.

5. APP configuration extension
- Add the property to config.yaml
- Add the property to the Struct 'model.Config'
- Read the config
```
	config := c.MustGet("config").(model.Config)
```

## Screenshots
### Guest screens
- Login:
![ScreenShot](/documentation/login.png)
- Register: 
![ScreenShot](/documentation/register.png)
### Administration back office
- Menu:
![ScreenShot](/documentation/menu.png)
- Casbin groups manager:
![ScreenShot](/documentation/groups_manager.png)
- Users manager:
![ScreenShot](/documentation/user_manager.png)
- User details:
![ScreenShot](/documentation/user_administration.png)
- Casbin URL resource manager:
![ScreenShot](/documentation/url_access_manager.png)


#### Recommendations
- Create separate Domain Struct and Form Struct