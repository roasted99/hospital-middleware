package models

import "time"

type Staff struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-"`
	Hospital string		`json:"hospital" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type StaffCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

type StaffLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	StaffID int `json:"staff_id"`
	Username string `json:"username"`
	Hospital string `json:"hospital"`
}

