// handlers.user.go

package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"database/sql"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/models"
	"github.com/best-nazar/web-app/repository"
	"github.com/best-nazar/web-app/security"
	"github.com/gin-gonic/gin"
)

func showLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Login",
	}, "login.html")
}

func performLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := security.ComputeHmac256(c.PostForm("password"))

    var sameSiteCookie http.SameSite;
	user, recNum := repository.GetUserByUsername(username)
	// Check if the username/password combination is valid
	if models.IsPasswordValid(user, password) && recNum > 0 {
		// If the username/password is valid set the token in a cookie
		token := helpers.GenerateSessionToken(strconv.FormatUint(uint64(user.ID), 10))

		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.SetSameSite(sameSiteCookie)

		render(c, gin.H{
			"title": "Successful Login",
			"payload": &user}, "login-successful.html")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
	}
}

func logout(c *gin.Context) {
    var sameSiteCookie http.SameSite;

	// Clear the cookie
	c.SetCookie("token", "", -1, "", "", false, true)
	c.SetSameSite(sameSiteCookie)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func showRegistrationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	render(c, gin.H{
		"title": "Register"}, "register.html")
}

func register(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := security.ComputeHmac256(c.PostForm("password"))
	name := c.PostForm("name")
	birthday := helpers.StringToTimestamp(c.PostForm("birthday"))

    var sameSiteCookie http.SameSite;

	if u, err := registerNewUser(name, username, password, birthday); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := helpers.GenerateSessionToken(strconv.FormatUint(uint64(u.ID), 10))

		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.SetSameSite(sameSiteCookie)

		render(c, gin.H{
			"title": "Successful registration & Login",
			"payload": &u}, "login-successful.html")

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})
	}
}

// Register a new user with the given username and password
func registerNewUser(name, username, password string, birthday int64) (*models.User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("the password can't be empty")
	} else if _, r :=repository.GetUserByUsername(username); r > 0 {
		return nil, errors.New("the username isn't available")
	}

	user := models.User{
		Name: name,
		Birthday: sql.NullInt64{Int64: birthday, Valid: true},
		Username: username,
		Password: password,
	}

	repository.AddNewUser(&user)

	return &user, nil
}