package models

import "time"

type User struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required"`
	PhoneNumber *string `json:"phone_number" binding:"required"`
	Type        string  `json:"type" binding:"required"`
}

type UpdateUser struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	PhoneNumber *string `json:"phone_number" binding:"required"`
	Type        string  `json:"type" binding:"required"`
}

type GetUser struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber *string   `json:"phone_number"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type GetAllUsers struct {
	Users []*GetUser `json:"users"`
	Count int64   `json:"count"`
}
