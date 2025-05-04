package middleware

import (
	"net/http"
	"strings"
	"context"

	"github.com/roasted99/hospital-middleware/internal/services"
	"github.com/roasted99/hospital-middleware/internal/utils"
)

type StaffContext string

const StaffKey StaffContext = "staff"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		staff, err := services.ValidateToken(token)
		if err != nil {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), StaffKey, staff)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}