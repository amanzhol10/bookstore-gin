package handlers

import (
	"bookstore/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var books = []models.Book{}
var bookID = 1

func GetBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if book.Title == "" || book.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	book.ID = bookID
	bookID++

	books = append(books, book)

	c.JSON(http.StatusCreated, book)
}

func GetBookByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for _, b := range books {
		if b.ID == id {
			c.JSON(http.StatusOK, b)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

func UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, b := range books {
		if b.ID == id {
			var updated models.Book

			if err := c.ShouldBindJSON(&updated); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			updated.ID = id
			books[i] = updated

			c.JSON(http.StatusOK, updated)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}

func DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}
