package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/roasted99/hospital-middleware/internal/api/middleware"
	"github.com/roasted99/hospital-middleware/internal/models"
	"github.com/roasted99/hospital-middleware/internal/services"
	"github.com/roasted99/hospital-middleware/internal/utils"
)

func SearchPatient(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		staffCtx := r.Context().Value(middleware.StaffKey)
		if staffCtx == nil {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		staff := staffCtx.(*models.Staff)

		query := models.PatientSearchRequest{
			NationalID:  r.URL.Query().Get("national_id"),
			PassportID:  r.URL.Query().Get("passport_id"),
			FirstName:   r.URL.Query().Get("first_name"),
			MiddleName:  r.URL.Query().Get("middle_name"),
			LastName:    r.URL.Query().Get("last_name"),
			DateOfBirth: r.URL.Query().Get("date_of_birth"),
			PhoneNumber: r.URL.Query().Get("phone_number"),
			Email:       r.URL.Query().Get("email"),
		}

		if strings.EqualFold(staff.Hospital, "Hospital A") {
			if query.NationalID != "" || query.PassportID != "" {
				searchID := query.NationalID
				if searchID == "" {
					searchID = query.PassportID
				}

				client := services.NewHospitalAClient()

				patient, err := client.SearchPatient(searchID)
				if err == nil {
					utils.ResponseWithSuccess(w, http.StatusOK, patient)
					return
				}
			}

			var queryArgs []interface{}
			var conditions []string
			var counter int = 1

			sqlQuery := "SELECT * FROM patient WHERE hospital = $1"
			queryArgs = append(queryArgs, staff.Hospital)
			counter++

			if query.NationalID != "" {
				conditions = append(conditions, "national_id = $"+strconv.Itoa(counter))
				queryArgs = append(queryArgs, query.NationalID)
				counter++
			}

			if query.PassportID != "" {
				conditions = append(conditions, "passport_id = $"+strconv.Itoa(counter))
				queryArgs = append(queryArgs, query.PassportID)
				counter++
			}

			if query.FirstName != "" {
				conditions = append(conditions, "first_name_en ILIKE $"+strconv.Itoa(counter)+" OR first_name_th ILIKE $"+strconv.Itoa(counter+1))
				queryArgs = append(queryArgs, "%"+query.FirstName+"%")
				queryArgs = append(queryArgs, "%"+query.FirstName+"%")
				counter +=2
			}

			if query.MiddleName != "" {
				conditions = append(conditions, "middle_name_en ILIKE $"+strconv.Itoa(counter)+" OR middle_name_th ILIKE $"+strconv.Itoa(counter+1))
				queryArgs = append(queryArgs, "%"+query.MiddleName+"%")
				queryArgs = append(queryArgs, "%"+query.MiddleName+"%")
				counter += 2
			}

			if query.LastName != "" {
				conditions = append(conditions, "last_name_en ILIKE $"+strconv.Itoa(counter)+" OR last_name_th ILIKE $"+strconv.Itoa(counter+1))
				queryArgs = append(queryArgs, "%"+query.LastName+"%")
				queryArgs = append(queryArgs, "%"+query.LastName+"%")
				counter += 2
			}

			if query.DateOfBirth != "" {
				conditions = append(conditions, "date_of_birth::text LIKE $"+strconv.Itoa(counter))
				queryArgs = append(queryArgs, "%"+query.DateOfBirth+"%")
				counter++
			}

			if query.PhoneNumber != "" {
				conditions = append(conditions, "phone_number ILIKE $"+strconv.Itoa(counter))
				queryArgs = append(queryArgs, "%"+query.PhoneNumber+"%")
				counter++
			}

			if query.Email != "" {
				conditions = append(conditions, "email ILIKE $"+strconv.Itoa(counter))
				queryArgs = append(queryArgs, "%"+query.Email+"%")
				counter++
			}

			if len(conditions) > 0 {
				sqlQuery += " AND " + strings.Join(conditions, " AND ")
			}
			fmt.Println(sqlQuery)
			fmt.Println(queryArgs)

			rows, err := db.Query(sqlQuery, queryArgs...)
			if err != nil {
				fmt.Println(err)
				utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to search patient")
				return
			}
			defer rows.Close()

			var patients []models.Patient
			for rows.Next() {
				var p models.Patient
				var middleNameTH, middleNameEN, nationalID, passportID sql.NullString
				err := rows.Scan(&p.ID, &p.FirstNameTH, &middleNameTH, &p.LastNameTH, &p.FirstNameEN, &middleNameEN, &p.LastNameEN, &p.DateOfBirth, &p.PatientHN, &nationalID, &passportID, &p.PhoneNumber, &p.Email, &p.Gender, &p.Hospital, &p.CreatedAt, &p.UpdatedAt)

				if middleNameTH.Valid {
					p.MiddleNameTH = middleNameTH.String
				}

				if middleNameEN.Valid {
					p.MiddleNameEN = middleNameEN.String
				}

				if nationalID.Valid {
					p.NationalID = nationalID.String
				}

				if passportID.Valid {
					p.PassportID = passportID.String
				}

				if err != nil {
					utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to scan patient")
					fmt.Println(err)
					return
				}
				patients = append(patients, p)
			}
			fmt.Println(patients)
			if len(patients) == 0 {
				utils.ResponseWithError(w, http.StatusNotFound, "No patient found")
				return
			}
			utils.ResponseWithSuccess(w, http.StatusOK, patients)
		} else {
			utils.ResponseWithError(w, http.StatusBadRequest, staff.Hospital+" is not supported yet")
		}
	}

}
