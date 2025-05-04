package models

import (
	"time"
)

type Patient struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	FirstNameTH string    `json:"first_name_th"`
	MiddleNameTH string    `json:"middle_name_th"`
	LastNameTH string    `json:"last_name_th"`
	FirstNameEN string    `json:"first_name_en"`
	MiddleNameEN string    `json:"middle_name_en"`
	LastNameEN string    `json:"last_name_en"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PatientHN string    `json:"patient_hn"`
	NationalID string    `json:"national_id"`
	PassportID string    `json:"passport_id"`
	PhoneNumber string    `json:"phone_number"`
	Email string    `json:"email"`
	Gender string
	Hospital string `json:"hospital"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PatientSearchRequest struct {
	NationalID string `json:"national_id"`
	PassportID string `json:"passport_id"`
	FirstName string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	PhoneNumber string `json:"phone_number"`
	Email string `json:"email"`
}