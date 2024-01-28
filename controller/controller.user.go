// controller.user.go

package controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/best-nazar/web-app/security"
	"github.com/gin-gonic/gin"
)

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
	password := security.ComputeHmac256(c.PostForm("password"))

	user, recNum := repository.GetUserByUsername(username)
	// Check if the username/password combination is valid
	if model.IsPasswordValid(*user, password) && recNum > 0 {
		// If the username/password is valid set the token in a cookie
		saveAuthToken(c, user)

		Render(c, gin.H{
			"title":   "Successful Login",
			"payload": &user}, "login-successful.html", http.StatusOK)
	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		er := errors.New("Invalid credentials provided")
		c.Error(er)
		Render(c, gin.H{"errors": c.Errors}, "login.html", http.StatusBadRequest)
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
	username := c.PostForm("username")

	if c.PostForm("password") != c.PostForm("password_repeat") {
		er := errors.New("Provided passwords do not match")
		c.Error(er)
		Render(c, gin.H{"errors": c.Errors}, "register.html", http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	birthday := helpers.StringToTimestamp(c.PostForm("birthday"))
	config, exist := c.Get("config")
	role := config.(model.Config)

	if !exist {
		panic("The Key 'default-casbin-group' is not found in config.yaml")
	}

	password := security.ComputeHmac256(c.PostForm("password"))

	if u, err := registerNewUser(name, username, password, role.DefaultCasbinGroup, birthday); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		saveAuthToken(c, u)

		Render(c, gin.H{
			"title":   "Successful registration & Login",
			"payload": &u}, "login-successful.html", http.StatusOK)

	} else {
		// If the username/password combination is invalid,
		// show the error message on the Register page
		c.Error(err)
		Render(c, gin.H{"errors": c.Errors}, "register.html", http.StatusBadRequest)
	}
}

// Register a new user with the given username and password
func registerNewUser(name, username, password, role string, birthday int64) (*model.User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("the password can't be empty")
	} else if _, r := repository.GetUserByUsername(username); r > 0 {
		return nil, errors.New("the username isn't available")
	}

	user := model.User{
		Name:     name,
		Birthday: sql.NullInt64{Int64: birthday, Valid: true},
		Username: username,
		Password: password,
	}

	repository.AddNewUser(&user)
	repository.AddCasbinUserRole(username, role)

	return &user, nil
}

// Save Auth data and token for UI
func saveAuthToken(c *gin.Context, user *model.User) {
	token := helpers.GenerateSessionToken(strconv.FormatUint(uint64(user.ID), 10))
	c.SetCookie("token", token, 3600, "", "", false, true)
	c.Set("is_logged_in", true)
	var sameSiteCookie http.SameSite
	c.SetSameSite(sameSiteCookie)
	c.Request.SetBasicAuth(user.Username, user.Password)
	c.SetCookie("auth", c.Request.Header.Get("Authorization"), 3600, "", "", false, true)
}
