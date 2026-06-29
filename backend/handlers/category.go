package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"careerhubv2-backend/config"
	"careerhubv2-backend/models"
)

// GET /api/v1/categories
func GetCategories(c *gin.Context) {
	limitStr := c.Query("limit")
	sortBy := c.Query("sort")

	// Validate sort param
	if sortBy != "" && sortBy != "job_count" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'sort' supports only 'job_count'",
			"code":  "INVALID_QUERY_PARAMETER",
		})
		return
	}

	// Validate limit param
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Query parameter 'limit' must be a positive integer",
				"code":  "INVALID_QUERY_PARAMETER",
			})
			return
		}
	}

	query := `
		SELECT
			c.CategoryID,
			c.CategoryName,
			ISNULL(c.IconName, '') AS IconName,
			COUNT(CASE
				WHEN j.IsActive = 1
				AND (j.DeadlineAt IS NULL OR j.DeadlineAt > GETDATE())
				THEN 1
			END) AS ActiveJobCount
		FROM Categories c
		LEFT JOIN Jobs j ON j.CategoryID = c.CategoryID
		GROUP BY c.CategoryID, c.CategoryName, c.IconName`

	if sortBy == "job_count" {
		query += " ORDER BY ActiveJobCount DESC"
	} else {
		query += " ORDER BY c.CategoryID ASC"
	}

	if limit > 0 {
		query += " OFFSET 0 ROWS FETCH NEXT " + strconv.Itoa(limit) + " ROWS ONLY"
	}

	rows, err := config.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve categories",
			"code":  "DATABASE_ERROR",
		})
		return
	}
	defer rows.Close()

	categories := []models.Category{}
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.CategoryID, &cat.CategoryName, &cat.IconName, &cat.ActiveJobCount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to parse categories",
				"code":  "DATABASE_ERROR",
			})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}
