package models

// JobSummary is returned by GET /jobs (list view)
type JobSummary struct {
	JobID        int      `json:"id"`
	JobTitle     string   `json:"job_title"`
	CompanyName  string   `json:"company_name"`
	CategoryID   int      `json:"category_id"`
	CategoryName string   `json:"category_name"`
	JobType      string   `json:"job_type"`
	State        *string  `json:"state"`
	City         *string  `json:"city"`
	SalaryMin    *float64 `json:"salary_min"`
	SalaryMax    *float64 `json:"salary_max"`
	Deadline     *string  `json:"deadline"`
	IsActive     bool     `json:"is_active"`
	CreatedAt    string   `json:"created_at"`
}

// JobDetail is returned by GET /jobs/:id (detail view, all tabs)
type JobDetail struct {
	JobID                 int      `json:"id"`
	JobTitle              string   `json:"job_title"`
	CompanyName           string   `json:"company_name"`
	CategoryID            int      `json:"category_id"`
	CategoryName          string   `json:"category_name"`
	JobType               string   `json:"job_type"`
	State                 *string  `json:"state"`
	City                  *string  `json:"city"`
	SalaryMin             *float64 `json:"salary_min"`
	SalaryMax             *float64 `json:"salary_max"`
	Deadline              *string  `json:"deadline"`
	JobDescription        *string  `json:"job_description"`
	Responsibilities      *string  `json:"responsibilities"`
	Requirements          *string  `json:"requirements"`
	AdditionalInformation *string  `json:"additional_information"`
	HowToApply            *string  `json:"how_to_apply"`
	IsActive              bool     `json:"is_active"`
	CreatedAt             string   `json:"created_at"`
	UpdatedAt             string   `json:"updated_at"`
}

type Pagination struct {
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

type JobsResponse struct {
	Data       []JobSummary `json:"data"`
	Pagination Pagination   `json:"pagination"`
}
