package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServiceImpl struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretKey string, issuer string) JWTService {
	return &JWTServiceImpl{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

func (j *JWTServiceImpl) GenerateToken(userId string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"iss":     j.issuer,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTServiceImpl) ValidateToken(tokenString string) (string, error) {
	if strings.TrimSpace(tokenString) == "" {
		return "", errors.New("token is required")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return "", err
	}
	if token == nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", errors.New("invalid token user id")
	}

	return userID, nil
}
