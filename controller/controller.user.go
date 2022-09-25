// controller.user.go

package controller

import (
	"net/http"
	"strconv"
	"strings"
	"errors"
	"database/sql"

	"github.com/best-nazar/web-app/helpers"
	"github.com/best-nazar/web-app/model"
	"github.com/best-nazar/web-app/repository"
	"github.com/best-nazar/web-app/security"
	"github.com/gin-gonic/gin"
	"github.com/best-nazar/web-app/errorSrc"
)

func ShowLoginPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title": "Login",
		"payload": "Login page",
	}, "login.html", http.StatusOK)
}

func PerformLogin(c *gin.Context) {
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := security.ComputeHmac256(c.PostForm("password"))

    var sameSiteCookie http.SameSite;
	user, recNum := repository.GetUserByUsername(username)
	// Check if the username/password combination is valid
	if model.IsPasswordValid(user, password) && recNum > 0 {
		// If the username/password is valid set the token in a cookie
		token := helpers.GenerateSessionToken(strconv.FormatUint(uint64(user.ID), 10))

		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.SetSameSite(sameSiteCookie)
		c.Request.SetBasicAuth(user.Username, user.Password)
		c.SetCookie("auth", c.Request.Header.Get("Authorization"), 3600, "", "", false, true)

		Render(c, gin.H{
			"title": "Successful Login",
			"payload": &user}, "login-successful.html", http.StatusOK)

	} else {
		// If the username/password combination is invalid,
		// show the error message on the login page
		err := errorSrc.ErrorView{"Login Failed", "Invalid credentials provided"}
		Render(c, gin.H{"error":  err},"login.html", http.StatusBadRequest)
	}
}

func Logout(c *gin.Context) {
    var sameSiteCookie http.SameSite;

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
	// Obtain the POSTed username and password values
	username := c.PostForm("username")
	password := security.ComputeHmac256(c.PostForm("password"))

	if err := helpers.ValidateUserPassword(password, c.PostForm("password_repeat")); err != nil {
		er := errorSrc.ErrorView{"Passwords", err.Error()}
		Render(c, gin.H{"error":  er},"register.html", http.StatusBadRequest)
		return
	}

	name := c.PostForm("name")
	birthday := helpers.StringToTimestamp(c.PostForm("birthday"))

    var sameSiteCookie http.SameSite;

	if u, err := registerNewUser(name, username, password, birthday); err == nil {
		// If the user is created, set the token in a cookie and log the user in
		token := helpers.GenerateSessionToken(strconv.FormatUint(uint64(u.ID), 10))

		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		c.SetSameSite(sameSiteCookie)

		Render(c, gin.H{
			"title": "Successful registration & Login",
			"payload": &u}, "login-successful.html", http.StatusOK)

	} else {
		// If the username/password combination is invalid,
		// show the error message on the Register page
		er := errorSrc.ErrorView{"Registration Failed", err.Error()}
		Render(c, gin.H{"error":  er},"register.html", http.StatusBadRequest)
	}
}

// Register a new user with the given username and password
func registerNewUser(name, username, password string, birthday int64) (*model.User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("the password can't be empty")
	} else if _, r :=repository.GetUserByUsername(username); r > 0 {
		return nil, errors.New("the username isn't available")
	}

	user := model.User{
		Name: name,
		Birthday: sql.NullInt64{Int64: birthday, Valid: true},
		Username: username,
		Password: password,
	}

	repository.AddNewUser(&user)

	return &user, nil
}