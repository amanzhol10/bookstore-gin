package handlers

import (
	"bookstore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var authors = []models.Author{}
var authorID = 1

func GetAuthors(c *gin.Context) {
	c.JSON(http.StatusOK, authors)
}

func CreateAuthor(c *gin.Context) {
	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if author.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name required"})
		return
	}

	author.ID = authorID
	authorID++

	authors = append(authors, author)

	c.JSON(http.StatusCreated, author)
}
