package handlers

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"careerhubv2-backend/config"
	"careerhubv2-backend/models"
)

// GET /api/v1/jobs
func GetJobs(c *gin.Context) {
	// Parse and validate query params
	page, err := parsePositiveInt(c.DefaultQuery("page", "1"), "page")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "INVALID_QUERY_PARAMETER"})
		return
	}

	limit, err := parsePositiveInt(c.DefaultQuery("limit", "10"), "limit")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "INVALID_QUERY_PARAMETER"})
		return
	}
	if limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameter 'limit' must not exceed 100",
			"code":  "INVALID_QUERY_PARAMETER",
		})
		return
	}

	search := strings.TrimSpace(c.Query("search"))

	var categoryID, jobTypeID *int
	if v := c.Query("category_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Query parameter 'category_id' must be a valid positive integer",
				"code":  "INVALID_QUERY_PARAMETER",
			})
			return
		}
		categoryID = &id
	}
	if v := c.Query("job_type_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Query parameter 'job_type_id' must be a valid positive integer",
				"code":  "INVALID_QUERY_PARAMETER",
			})
			return
		}
		jobTypeID = &id
	}

	// Build WHERE clause dynamically
	where := []string{
		"j.IsActive = 1",
		"(j.DeadlineAt IS NULL OR j.DeadlineAt > GETDATE())",
	}
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		where = append(where, fmt.Sprintf(
			"(j.JobTitle LIKE @p%d OR j.CompanyName LIKE @p%d OR c.CategoryName LIKE @p%d)",
			argIdx, argIdx, argIdx,
		))
		args = append(args, "%"+search+"%")
		argIdx++
	}
	if categoryID != nil {
		where = append(where, fmt.Sprintf("j.CategoryID = @p%d", argIdx))
		args = append(args, *categoryID)
		argIdx++
	}
	if jobTypeID != nil {
		where = append(where, fmt.Sprintf("j.JobTypeID = @p%d", argIdx))
		args = append(args, *jobTypeID)
		argIdx++
	}

	whereClause := "WHERE " + strings.Join(where, " AND ")

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM Jobs j
		INNER JOIN Categories c ON c.CategoryID = j.CategoryID
		%s`, whereClause)

	var totalRecords int
	if err := config.DB.QueryRow(countQuery, args...).Scan(&totalRecords); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count jobs", "code": "DATABASE_ERROR"})
		return
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}
	offset := (page - 1) * limit

	// Fetch page
	dataQuery := fmt.Sprintf(`
		SELECT
			j.JobID, j.JobTitle, j.CompanyName,
			j.CategoryID, c.CategoryName,
			jt.TypeName,
			j.State, j.City,
			j.SalaryMin, j.SalaryMax,
			CONVERT(VARCHAR, j.DeadlineAt, 127),
			j.IsActive,
			CONVERT(VARCHAR, j.CreatedAt, 127)
		FROM Jobs j
		INNER JOIN Categories c  ON c.CategoryID  = j.CategoryID
		INNER JOIN JobTypes   jt ON jt.JobTypeID  = j.JobTypeID
		%s
		ORDER BY j.CreatedAt DESC
		OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY`,
		whereClause, argIdx, argIdx+1)

	args = append(args, offset, limit)

	rows, err := config.DB.Query(dataQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs", "code": "DATABASE_ERROR"})
		return
	}
	defer rows.Close()

	jobs := []models.JobSummary{}
	for rows.Next() {
		var j models.JobSummary
		var isActiveBit int
		if err := rows.Scan(
			&j.JobID, &j.JobTitle, &j.CompanyName,
			&j.CategoryID, &j.CategoryName,
			&j.JobType,
			&j.State, &j.City,
			&j.SalaryMin, &j.SalaryMax,
			&j.Deadline,
			&isActiveBit,
			&j.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse jobs", "code": "DATABASE_ERROR"})
			return
		}
		j.IsActive = isActiveBit == 1
		jobs = append(jobs, j)
	}

	c.JSON(http.StatusOK, models.JobsResponse{
		Data: jobs,
		Pagination: models.Pagination{
			Page:         page,
			Limit:        limit,
			TotalRecords: totalRecords,
			TotalPages:   totalPages,
		},
	})
}

// GET /api/v1/jobs/:id
func GetJobByID(c *gin.Context) {
	idStr := c.Param("id")
	jobID, err := strconv.Atoi(idStr)
	if err != nil || jobID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Path parameter 'id' must be a valid positive integer",
			"code":  "INVALID_PATH_PARAMETER",
		})
		return
	}

	var j models.JobDetail
	var isActiveBit int
	err = config.DB.QueryRow(`
		SELECT
			j.JobID, j.JobTitle, j.CompanyName,
			j.CategoryID, c.CategoryName,
			jt.TypeName,
			j.State, j.City,
			j.SalaryMin, j.SalaryMax,
			CONVERT(VARCHAR, j.DeadlineAt, 127),
			j.JobDescription, j.Responsibilities, j.Requirements,
			j.AdditionalInformation, j.HowToApply,
			j.IsActive,
			CONVERT(VARCHAR, j.CreatedAt, 127),
			CONVERT(VARCHAR, j.UpdatedAt, 127)
		FROM Jobs j
		INNER JOIN Categories c  ON c.CategoryID = j.CategoryID
		INNER JOIN JobTypes   jt ON jt.JobTypeID = j.JobTypeID
		WHERE j.JobID = @p1`,
		jobID,
	).Scan(
		&j.JobID, &j.JobTitle, &j.CompanyName,
		&j.CategoryID, &j.CategoryName,
		&j.JobType,
		&j.State, &j.City,
		&j.SalaryMin, &j.SalaryMax,
		&j.Deadline,
		&j.JobDescription, &j.Responsibilities, &j.Requirements,
		&j.AdditionalInformation, &j.HowToApply,
		&isActiveBit,
		&j.CreatedAt, &j.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Job listing with ID '%d' does not exist", jobID),
			"code":  "RESOURCE_NOT_FOUND",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve job details",
			"code":  "DATABASE_ERROR",
		})
		return
	}

	j.IsActive = isActiveBit == 1
	c.JSON(http.StatusOK, j)
}

func parsePositiveInt(val, paramName string) (int, error) {
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("Query parameter '%s' must be a positive integer", paramName)
	}
	return n, nil
}
