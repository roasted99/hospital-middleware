package handlers_test

import (
	"bytes"
	// "database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/roasted99/hospital-middleware/internal/api/handlers"
	"github.com/roasted99/hospital-middleware/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateStaff(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name           string
		requestBody    models.StaffCreateRequest
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful staff creation",
			requestBody: models.StaffCreateRequest{
				Username: "testuser",
				Password: "password123",
				Hospital: "Test Hospital",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO staff").
					WithArgs("testuser", sqlmock.AnyArg(), "Test Hospital", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"token":    "valid_jwt_token",
					"staff_id": 1,
					"username": "testuser",
					"hospital": "Test Hospital",
				},
				"status":  "Created",
				"message": "Success",
			},
		},
		{
			name: "Failed staff creation due to missing fields",
			requestBody: models.StaffCreateRequest{
				Username: "",
				Password: "password123",
				Hospital: "Test Hospital",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"status":  "Bad Request",
				"message": "Username, password, and hospital are required",
			},
		},
		{
			name: "Failed staff creation due to database error",
			requestBody: models.StaffCreateRequest{
				Username: "testuser",
				Password: "password123",
				Hospital: "Test Hospital",
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO staff").
					WithArgs("testuser", sqlmock.AnyArg(), "Test Hospital", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"status":  "Internal Server Error",
				"message": "Database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler := handlers.CreateStaff(db)
			handler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			for key := range tt.expectedBody {
				// assert.Equal(t, expectedValue, response[key])
				assert.Contains(t, response, key)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestLoginStaff(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)

	tests := []struct {
		name           string
		requestBody    models.StaffLoginRequest
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful staff login",
			requestBody: models.StaffLoginRequest{
				Username: "testuser",
				Password: "password123",
				Hospital: "Test Hospital",
			},
			mockSetup: func() {
					mock.ExpectQuery("SELECT id, username, password, hospital FROM staff WHERE username = \\$1 AND hospital = \\$2").
					WithArgs("testuser", "Test Hospital").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "hospital"}).AddRow(1, "testuser", string(hashedPassword), "Test Hospital"))
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"token":    "valid_jwt_token",
					"staff_id": float64(1),
					"username": "testuser",
				},
				"status":  "OK",
				"message": "Success",
			},
		},
		{
			name: "Failed staff login due to invalid credentials",
			requestBody: models.StaffLoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
				Hospital: "Test Hospital",
			},
			mockSetup: func() {
				mock.ExpectQuery("SELECT id, username, password, hospital FROM staff WHERE username = \\$1 AND hospital = \\$2").
					WithArgs("testuser", "Test Hospital").
					WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "hospital"}).AddRow(1, "testuser", "hashed_password","Test Hospital"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"status":  "Unauthorized",
				"message": "Invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler := handlers.LoginStaff(db)
			handler(w, req)

			res := w.Result()
			assert.Equal(t, tt.expectedStatus, res.StatusCode)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			for key := range tt.expectedBody {
				assert.Contains(t, response, key)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
