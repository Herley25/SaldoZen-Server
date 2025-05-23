package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func genToken(userID string, secret []byte, minutes int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(minutes) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateAccess(userID string) (string, error) {
	minutes, _ := strconv.Atoi(os.Getenv("ACCESS_EXPIRED_MINUTES"))
	if minutes == 0 {
		minutes = 15
	}
	return genToken(userID, []byte(os.Getenv("JWT_SECRET")), minutes)
}

func GenerateRefresh(userID string) (string, error) {
	days, _ := strconv.Atoi(os.Getenv("REFRESH_EXPIRED_DAYS"))
	if days == 0 {
		days = 7
	}
	minutes := days * 24 * 60
	return genToken(userID, []byte(os.Getenv("JWT_REFRESH_SECRET")), minutes)
}
