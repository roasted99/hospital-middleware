package handlers_test

import (
	"context"
	"encoding/json"
	// "fmt"
	"net/http"
	"net/http/httptest"
	"time"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/roasted99/hospital-middleware/internal/api/handlers"
	"github.com/roasted99/hospital-middleware/internal/api/middleware"
	"github.com/roasted99/hospital-middleware/internal/models"
)

func createAuthenticatedRequest(method, url string, staff *models.Staff) *http.Request {
	req, _ := http.NewRequest(method, url, nil)
	ctx := context.WithValue(req.Context(), middleware.StaffKey, staff)
	return req.WithContext(ctx)
}

func TestSearchPatient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	tests := []struct {
		name           string
		staff          *models.Staff
		url            string
		mockSetup      func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Valid request with National ID",
			staff: &models.Staff{
				Hospital: "Hospital A",
				Username: "staff1",
				ID:       1,
			},
			url: "/patient/search?national_id=1234567890123",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM patient WHERE hospital = \\$1 AND national_id = \\$2").
					WithArgs("Hospital A", "1234567890123").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name_th", "middle_name_th", "last_name_th", "first_name_en", "middle_name_en", "last_name_en", "date_of_birth", "patient_hn", "national_id", "passport_id", "phone_number", "email", "gender", "hospital", "created_at", "updated_at"}).
						AddRow(1, "ทดสอบ", "กลาง", "สุดท้าย", "Test", "Middle", "Last", time.Now(), "HN123456", "1234567890123", "", "0123456789", "test@email.com", "M", "Hospital A", time.Now(), time.Now()))
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Success",
				"data": map[string]interface{}{
					"id":             1,
					"first_name_th":  "ทดสอบ",
					"middle_name_th": "กลาง",
					"last_name_th":   "สุดท้าย",
					"first_name_en":  "Test",
					"middle_name_en": "Middle",
					"last_name_en":   "Last",
					"date_of_birth":  "2000-01-01",
					"patient_hn":     "HN123456",
					"national_id":    "1234567890123",
					"passport_id":    "",
					"phone_number":   "0123456789",
					"email":          "test@email.com",
					"gender":         "M",
					"hospital":       "Hospital A",
					"created_at":     "2023-01-01",
					"updated_at":     "2023-01-01",
				},
			},
		},
		{
			name: "Search by multiple parameters",
			staff: &models.Staff{
				Hospital: "Hospital A",
				Username: "staff1",
				ID:       1,
			},
			url: "/patient/search?first_name=Test&last_name=Last",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM patient WHERE hospital = \\$1 AND first_name_en ILIKE \\$2 OR first_name_th ILIKE \\$3 AND last_name_en ILIKE \\$4 OR last_name_th ILIKE \\$5").
					WithArgs("Hospital A", "%"+"Test"+"%", "%"+"Test%", "%Last%", "%Last%").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name_th", "middle_name_th", "last_name_th", "first_name_en", "middle_name_en", "last_name_en", "date_of_birth", "patient_hn", "national_id", "passport_id", "phone_number", "email", "gender", "hospital", "created_at", "updated_at"}).
						AddRow(1, "ทดสอบ", "กลาง", "สุดท้าย", "Test", "Middle", "Last", time.Now(), "HN123456", "1234567890123", "", "0123456789", "test@email.com", "M", "Hospital A", time.Now(), time.Now()))
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":  "OK",
				"message": "Success",
				"data": map[string]interface{}{
					"id":             1,
					"first_name_th":  "ทดสอบ",
					"middle_name_th": "กลาง",
					"last_name_th":   "สุดท้าย",
					"first_name_en":  "Test",
					"middle_name_en": "Middle",
					"last_name_en":   "Last",
					"date_of_birth":  "2000-01-01",
					"patient_hn":     "HN123456",
					"national_id":    "1234567890123",
					"passport_id":    "",
					"phone_number":   "0123456789",
					"email":          "test@email.com",
					"gender":         "M",
					"hospital":       "Hospital A",
					"created_at":     "2023-01-01",
					"updated_at":     "2023-01-01",
				},
			},
		},
		{
			name: "No patient found",
			staff: &models.Staff{
				Hospital: "Hospital A",
				Username: "staff1",
				ID:       1,
			},
			url: "/patient/search?passport_id=12345678",
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM patient WHERE hospital = \\$1 AND passport_id = \\$2").
					WithArgs("Hospital A", "12345678").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name_th", "middle_name_th", "last_name_th", "first_name_en", "middle_name_en", "last_name_en", "date_of_birth", "patient_hn", "national_id", "passport_id", "phone_number", "email", "gender", "hospital"}))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"status":  "Not Found",
				"message": "No patient found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock database
			tt.mockSetup()
			// Create request
			req := createAuthenticatedRequest("GET", tt.url, tt.staff)

			rr := httptest.NewRecorder()

			handler := handlers.SearchPatient(db)

			handler(rr, req)

			// Check status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(rr.Body.Bytes(), &response)

				if response["status"] != tt.expectedBody["status"] || response["message"] != tt.expectedBody["message"] {
					t.Errorf("expected body %v, got %v", tt.expectedBody, response)
				}
				
				if data, ok := response["data"].(map[string]interface{}); ok {
					if data["id"] != tt.expectedBody["data"].(map[string]interface{})["id"] {
						t.Errorf("expected patient ID %v, got %v", tt.expectedBody["data"].(map[string]interface{})["id"], data["id"])
					}
				}
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Errorf("Unfulfilled expectations: %s", err)
				}
			}
		})
	}
}
