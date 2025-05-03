package models

import "time"

type Staff struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	HospitalID string		`json:"hospital_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type StaffCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	HospitalID string `json:"hospital_id" binding:"required"`
}

type StaffLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	HospitalID string `json:"hospital_id" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	StaffID int `json:"staff_id"`
	Username string `json:"username"`
	HospitalID string `json:"hospital_id"`
}

