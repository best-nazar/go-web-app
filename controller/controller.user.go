// controller.user.go

package controller

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/best-nazar/web-app/security"
	"github.com/gin-gonic/gin"
)

func ShowUserHomePage(c *gin.Context) {
	userLoggedIn, hasStatus := c.Get("is_logged_in")

	if hasStatus {
		isLogedIn := userLoggedIn.(bool)

		if isLogedIn {
			Render(c, gin.H{
				"title":   "Successful Login",
			}, "login-successful.html", http.StatusOK)
		}
	}

	c.AbortWithStatus(http.StatusForbidden)
}

func ShowLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Login",
		"payload": "Login page",
	}, "login.html", http.StatusOK)
}

func PerformLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := c.PostForm("password")

	user, recNum := repository.GetUserByUsername(username)
	// Check if the username/password combination is valid
	if user.IsPasswordValid(password) && recNum > 0 {
		// If the username/password is valid set the token in a cookie
		saveAuthToken(c, user)
		c.Redirect(http.StatusFound, "/u")
	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.Error(errors.New("user credentials|Invalid credentials provided"))
		Render(c, gin.H{"errors": helpers.Errors(c)}, "login.html", http.StatusBadRequest)
	}
}

func Logout(c *gin.Context) {
	var sameSiteCookie http.SameSite

	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)
	c.SetCookie("auth", "", -1, "", "", false, true)
	c.SetSameSite(sameSiteCookie)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func ShowRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{"title": "Register"}, "register.html", http.StatusOK)
}

func Register(c *gin.Context) {
	var user model.User

	e := c.ShouldBind(&user)

	if e != nil {
		c.Error(e)
	}

	if user.Password != c.PostForm("password_repeat") {
		er := errors.New("credentials|provided password does not match")
		c.Error(er)
	}

	if length := len(c.Errors); length > 0 {
		Render(c, gin.H{"errors": helpers.Errors(c)}, "register.html", http.StatusBadRequest)
		return
	}

	conf := c.MustGet("config").(model.Config)

	if slices.Contains(strings.Split(conf.UsernameRestrictedWords, ","), user.Username) {
		c.Error(errors.New("username|" + user.Username + " isn't available"))
		Render(c, gin.H{"errors": helpers.Errors(c)}, "register.html", http.StatusBadRequest)
		return
	}

	password := security.ComputeHmac256(c.PostForm("password"))

	if u, err := registerNewUser(&user, password, conf.DefaultCasbinGroup); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		saveAuthToken(c, u)

		Render(c, gin.H{
			"title": "Successful registration",
			"user":  &u,
			"admin": conf.ContactSupportEmail,
		}, "register-successful.html", http.StatusOK)

	} else {
		// If the username/password combination is invalid,
		// show the error message on the Register page
		c.Error(err)
		Render(c, gin.H{"errors": helpers.Errors(c)}, "register.html", http.StatusBadRequest)
	}
}

func UserLocked(c *gin.Context) {
	var user  *model.User
	var nRows int64
	var idModel = model.Id{}

	uerr := c.ShouldBindQuery(&idModel)

	if uerr == nil {
		user, nRows = repository.FindUserById(idModel.String())
	} else {
		c.Error(uerr)
	}

	if nRows == 0 {
		c.Error(errors.New("User ID|" +  idModel.String() + " not found"))
	} else {	
		conf := c.MustGet("config").(model.Config)

		Render(c, gin.H{
			"title": "Acount is locked ",
			"user":  &user,
			"admin": conf.ContactSupportEmail,
		}, "register-successful.html", http.StatusOK)

		return
	}

	Render(c, gin.H{
		"title": "Error",
		"description": "Account is locked",
		"errors" : helpers.Errors(c),
	}, "errors.html", http.StatusBadRequest)
}

// Register a new user with the given username and password
func registerNewUser(user *model.User, password string, role string) (*model.User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("password| can't be empty")
	} else if _, r := repository.GetUserByUsername(user.Username); r > 0 {
		return nil, errors.New("username| " + user.Username + " isn't available")
	}
	user.Password = password
	repository.AddNewUser(user)
	repository.AddCasbinUserRole(user.Username, role)

	return user, nil
}

// Save Auth data and token for UI
func saveAuthToken(c *gin.Context, user *model.User) {
	ipAddr := c.ClientIP()
	data := strconv.FormatUint(uint64(user.ID), 10)
	token := helpers.GenerateSessionToken(data, ipAddr)
	c.SetCookie("token", token, 3600, "", "", false, true)
	c.Set("is_logged_in", true)
	var sameSiteCookie http.SameSite
	c.SetSameSite(sameSiteCookie)
}
