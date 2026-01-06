package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func jwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateAccessToken(userID string, permissions []string) (string, error) {
	ttl := time.Minute * 15

	claims := jwt.MapClaims{
		"user_id":     userID,
		"permissions": permissions,
		"exp":         time.Now().Add(ttl).Unix(),
		"type":        "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

func GenerateRefreshToken(userID string, ttlDays int) (string, time.Time, error) {
	expires := time.Now().Add(time.Hour * 24 * time.Duration(ttlDays))

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expires.Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret())
	return signed, expires, err
}
