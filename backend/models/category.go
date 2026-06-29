package models

type Category struct {
	CategoryID     int    `json:"id"`
	CategoryName   string `json:"name"`
	IconName       string `json:"icon_name"`
	ActiveJobCount int    `json:"active_job_count"`
}
