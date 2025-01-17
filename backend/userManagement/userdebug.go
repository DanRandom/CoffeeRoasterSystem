package userManagement

import (
	"database/sql"
	"log"

	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func FetchUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User

		rows, err := db.Query("SELECT id, username FROM users")
		if err != nil {
			if err == sql.ErrNoRows {
				log.Fatal("The user doesn't exist in the database")
				c.JSON(http.StatusBadRequest, gin.H{"ERROR": "Users cannot be queried"})
			}
			return
		}

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.Id, &user.Username); err != nil {
				log.Printf("Failed to scan user: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
				return
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		query := "DELETE FROM users WHERE id = ?"
		res, err := db.Exec(query, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"title":          "RoastingRooster Coffee Inventory",
			"dbActionStatus": "SUCCESS",
		})
	}
}
