# Go gin app

This is the code from the article [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin).


1. An authorization middleware: https://github.com/gin-contrib/authz
2. Auth policy editor: https://casbin.org/editor/
3. Auth configuration: authz_model.conf & authz_policy.csv
4. In Casbin, an access control model is abstracted into a CONF file based on the PERM metamodel (Policy, Effect, Request, Matchers) https://github.com/casbin/casbin

5. Adding the policy rule: see automigratioon in 'dbconnection.go'
- Restart the App to reload rules. More examples at ``__authz_policy.csv``
6. Groups hierarchy
- guest
- - member
- - - admin

6.1. Admin user
6.2. To register a User with admin rights, change the configuration in config.yaml to the next 
```
default_casbin_group: admin
```
Then hit the '/u/register' URL and fill in the form. Then return the configuration in config.yaml to 'member'.

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
		"page-name.html",
		http.StatusOk})
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


8. User Activity Data collecting
Keeps the history of actions like: page opened, data added/updated/deleted
Turn on/off in config.yaml and check DB table 'user_activities'
```
user_activity_logging: true
```

9. DB Migration
5.1 Automigration: see "dbconnection.go"
Condition - if no tables created then auto run migration. In other words, drop all tables in schema, run the server, and hit any endpoint.

10. APP configuration
10.1. Add the propoerty to config.yaml
10.2. Add the property to the Struct 'model.Config'
10.3. Read the config
```
	config := c.MustGet("config").(model.Config)
```