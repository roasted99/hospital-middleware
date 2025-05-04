package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roasted99/hospital-middleware/internal/config"
	"github.com/roasted99/hospital-middleware/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	StaffID  int    `json:"staff_id"`
	Username string `json:"username"`
	Hospital string `json:"hospital"`
	jwt.RegisteredClaims
}

func GenerateJWT(staffID int, username, hospital string) (string, error) {
	claims := JWTClaims{
		StaffID:  staffID,
		Username: username,
		Hospital: hospital,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "hospital-middleware",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // Token expires in 72 hour
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWTSecret()))
}

func ValidateToken(tokenString string) (*models.Staff, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GetJWTSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	staff := &models.Staff{
		ID:      claims.StaffID,
		Username: claims.Username,
		Hospital: claims.Hospital,
	}

	return staff, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}