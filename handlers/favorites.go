package handlers

import (
	"bookstore/db"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddFavorite(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	found := false
	for _, b := range books {
		if b.ID == bookID {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	_, err = db.DB.Exec(
		`INSERT INTO favorite_books (user_id, book_id, created_at)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, book_id) DO NOTHING`,
		userID, bookID, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to favorites"})
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetInt("user_id")
	bookID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	res, err := db.DB.Exec(
		`DELETE FROM favorite_books WHERE user_id = $1 AND book_id = $2`,
		userID, bookID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "favorite not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func GetFavorites(c *gin.Context) {
	userID := c.GetInt("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var total int
	err := db.DB.QueryRow(
		`SELECT COUNT(*) FROM favorite_books WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := db.DB.Query(
		`SELECT book_id, created_at FROM favorite_books
		 WHERE user_id = $1
		 ORDER BY created_at DESC
		 LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	type FavoriteEntry struct {
		BookID    int       `json:"book_id"`
		CreatedAt time.Time `json:"created_at"`
	}

	favorites := []FavoriteEntry{}
	for rows.Next() {
		var f FavoriteEntry
		if err := rows.Scan(&f.BookID, &f.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		favorites = append(favorites, f)
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"data":      favorites,
	})
}
