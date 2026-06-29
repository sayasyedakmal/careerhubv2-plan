package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"careerhubv2-backend/config"
	"careerhubv2-backend/models"
)

// GET /api/v1/users/me
func GetMe(c *gin.Context) {
	userID := c.GetInt("userID")

	var u models.User
	err := config.DB.QueryRow(
		`SELECT UserID, Email, DisplayName, Phone, UserType, StudentID
		 FROM Users WHERE UserID = @p1`,
		userID,
	).Scan(&u.UserID, &u.Email, &u.DisplayName, &u.Phone, &u.UserType, &u.StudentID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Active user record not found in the database",
			"code":  "RESOURCE_NOT_FOUND",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, u)
}
