package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IssueJWT creates and signs a JWT for the given user ID.
func IssueJWT(userID int, secret string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
