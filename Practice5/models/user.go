package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	Birthdate time.Time `json:"birthdate"`
}

type PaginatedResponse struct {
	Data       []User `json:"data"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type FilterParams struct {
	ID        *int
	Name      *string
	Email     *string
	Gender    *string
	Birthdate *string // "YYYY-MM-DD"
	OrderBy   string  // "id", "name", "email", "gender", "birthdate"
	OrderDir  string  // "ASC" or "DESC"
	Page      int
	PageSize  int
}
