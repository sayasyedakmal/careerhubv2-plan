package models

type User struct {
	UserID             int     `json:"id"`
	Email              string  `json:"email"`
	DisplayName        string  `json:"display_name"`
	Phone              *string `json:"phone"`
	UserType           string  `json:"user_type"`
	StudentID          *string `json:"student_id"`
	MicrosoftObjectID  *string `json:"-"`
	PasswordHash       *string `json:"-"`
	RegistrationStatus string  `json:"-"`
}
