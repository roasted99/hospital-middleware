package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/roasted99/hospital-middleware/internal/models"
	"github.com/roasted99/hospital-middleware/internal/services"
	"github.com/roasted99/hospital-middleware/internal/utils"
)

func CreateStaff(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.StaffCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if request.Username == "" || request.Password == "" || request.Hospital == "" {
			utils.ResponseWithError(w, http.StatusBadRequest, "Username, password, and hospital are required")
			return
		}

		hashedPassword, err := services.HashPassword(request.Password)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}

		var staffID int
		err = db.QueryRow("INSERT INTO staff (username, password, hospital, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			request.Username, hashedPassword, request.Hospital, time.Now(), time.Now()).Scan(&staffID)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to create staff")
			} else {
				utils.ResponseWithError(w, http.StatusInternalServerError, "Database error")
			}
			return
		}

		token, err := services.GenerateJWT(staffID, request.Username, request.Hospital)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}
		
		utils.ResponseWithSuccess(w, http.StatusCreated, models.AuthResponse{
			Token:    token,
			StaffID:  staffID,
			Username: request.Username,
			Hospital: request.Hospital,
		})

	}
}

func LoginStaff(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.StaffLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if request.Username == "" || request.Password == "" || request.Hospital == "" {
			utils.ResponseWithError(w, http.StatusBadRequest, "Username, password, and hospital are required")
			return
		}

		var staff models.Staff
		err := db.QueryRow("SELECT id, username, password FROM staff WHERE username = $1 AND hospital = $2", request.Username, request.Hospital).Scan(&staff.ID, &staff.Username, &staff.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid credentials")
			} else {
				utils.ResponseWithError(w, http.StatusInternalServerError, "Database error")
			}
			return
		}

		if !services.CheckPasswordHash(request.Password, staff.Password) {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		token, err := services.GenerateJWT(staff.ID, staff.Username, request.Hospital)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, models.AuthResponse{
			Token:    token,
			StaffID:  staff.ID,
			Username: staff.Username,
			Hospital: request.Hospital,
		})
	}
}

