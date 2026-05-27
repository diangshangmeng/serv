package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"voucher-platform/config"
)

func GenerateToken(userID uint64, phone string, cityID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"phone":   phone,
		"city_id": cityID,
		"exp":     time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func GenerateAdminToken(adminID uint64, username string) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"username": username,
		"role":     "admin",
		"exp":      time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}