// controller.article.go

package controller

import (
	"net/http"
	"strconv"

	"github.com/best-nazar/web-app/model"
	"github.com/gin-gonic/gin"
)

func ShowIndexPage(c *gin.Context) {
	articles := model.GetAllArticles()

	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title":   "Home Page",
		"payload": articles}, "index.html", http.StatusOK)
}

func ShowArticleCreationPage(c *gin.Context) {
	// Call the render function with the name of the template to render
	Render(c, gin.H{
		"title": "Create New Article"}, "create-article.html", http.StatusOK)
}

func GetArticle(c *gin.Context) {
	// Check if the article ID is valid
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		// Check if the article exists
		if article, err := model.GetArticleByID(articleID); err == nil {
			// Call the render function with the title, article and the name of the
			// template
			Render(c, gin.H{
				"title":   article.Title,
				"payload": article}, "article.html", http.StatusOK)

		} else {
			// If the article is not found, abort with an error
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		// If an invalid article ID is specified in the URL, abort with an error
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func CreateArticle(c *gin.Context) {
	// Obtain the POSTed title and content values
	title := c.PostForm("title")
	content := c.PostForm("content")

	if a, err := model.CreateNewArticle(title, content); err == nil {
		// If the article is created successfully, show success message
		Render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "submission-successful.html", http.StatusOK)
	} else {
		// if there was an error while creating the article, abort with an error
		c.AbortWithStatus(http.StatusBadRequest)
	}
}