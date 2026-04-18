package main

import (
	"bookstore/db"
	"bookstore/handlers"
	"bookstore/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	r := gin.Default()

	r.GET("/books", handlers.GetBooks)
	r.POST("/books", handlers.CreateBook)
	r.GET("/books/:id", handlers.GetBookByID)
	r.PUT("/books/:id", handlers.UpdateBook)
	r.DELETE("/books/:id", handlers.DeleteBook)

	r.GET("/authors", handlers.GetAuthors)
	r.POST("/authors", handlers.CreateAuthor)
	r.GET("/categories", handlers.GetCategories)
	r.POST("/categories", handlers.CreateCategory)

	auth := r.Group("/", middleware.AuthRequired())
	{
		auth.GET("/books/favorites", handlers.GetFavorites)
		auth.PUT("/books/:id/favorites", handlers.AddFavorite)
		auth.DELETE("/books/:id/favorites", handlers.RemoveFavorite)
	}

	r.Run(":8000")
}
