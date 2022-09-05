# Go gin app

This is the code from the article [Building Go Web Applications and Microservices Using Gin](https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin).


1. An authorization middleware: https://github.com/gin-contrib/authz
2. Auth policy editor: https://casbin.org/editor/
3. Auth configuration: authz_model.conf & authz_policy.csv
4. In Casbin, an access control model is abstracted into a CONF file based on the PERM metamodel (Policy, Effect, Request, Matchers) https://github.com/casbin/casbin

5. Adding the policy rule:
- Run query: ```INSERT INTO casbin_rule (p_type, v0, v1, v2) VALUES ('p, anonymous, /article/view/*, GET);```
- Restart the App to reload rules.
6. Adding the group:
- Run query: ```INSERT INTO casbin_rule (p_type, v0, v1) VALUES ('g', 'test1', 'anonymous');```
- Restart the App to reload rules. More examples at ``__authz_policy.csv``