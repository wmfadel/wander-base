package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GernerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	secretKey := GetFromEnv("secretKey")
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) (int64, error) {
	parseedToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("wrong signing method")
			}
			return []byte(GetFromEnv("secretKey")), nil
		},
	)

	if err != nil {
		return 0, fmt.Errorf("jwt token parse failed %w", err)
	}

	if !parseedToken.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parseedToken.Claims.(jwt.MapClaims)

	if !ok {
		return 0, errors.New("invalid token")
	}

	// email, _ := claims["email"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, nil

}
