package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/roasted99/hospital-middleware/internal/config"
	"github.com/roasted99/hospital-middleware/internal/models"
)

type HospitalClient interface {
	SearchPatient(patientID string) (*models.Patient, error)
}

type HospitalAClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHospitalAClient() *HospitalAClient {
	return &HospitalAClient{
		BaseURL:    config.GetHospitalAURL(),
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type HospitalAResponse struct {
	FirstNameTH string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH string `json:"last_name_th"`
	FirstNameEN string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN string `json:"last_name_en"`
	DateOfBirth time.Time `json:"date_of_birth"`
	PatientHN string `json:"patient_hn"`
	NationalID string `json:"national_id"`
	PassportID string `json:"passport_id"`
	PhoneNumber string `json:"phone_number"`
	Email string `json:"email"`
	Gender string `json:"gender"`
}

func (c *HospitalAClient) SearchPatient(patientID string) (*models.Patient, error) {
	apiURL := fmt.Sprintf("%s/api/v1/patients/%s", c.BaseURL, url.PathEscape(patientID))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.GetJWTSecret())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to search patient: %s", resp.Status)
	}

	var hospitalAResponse HospitalAResponse
	if err := json.NewDecoder(resp.Body).Decode(&hospitalAResponse); err != nil {
		return nil, err
	}

	patient := &models.Patient{
		FirstNameTH: hospitalAResponse.FirstNameTH,
		MiddleNameTH: hospitalAResponse.MiddleNameTH,
		LastNameTH: hospitalAResponse.LastNameTH,
		FirstNameEN: hospitalAResponse.FirstNameEN,
		MiddleNameEN: hospitalAResponse.MiddleNameEN,
		LastNameEN: hospitalAResponse.LastNameEN,
		DateOfBirth: hospitalAResponse.DateOfBirth,
		PatientHN: hospitalAResponse.PatientHN,
		NationalID: hospitalAResponse.NationalID,
		PassportID: hospitalAResponse.PassportID,
		PhoneNumber: hospitalAResponse.PhoneNumber,
		Email: hospitalAResponse.Email,
		Gender: hospitalAResponse.Gender,
		Hospital: "Hospital A",
	}

	return patient, nil
}