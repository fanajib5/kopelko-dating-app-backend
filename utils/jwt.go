package utils

import (
	"fmt"
	"os"
	"time"

	model "kopelko-dating-app-backend/models"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

type Claims struct {
	UserID    uint
	Email     string
	IsPremium bool
	jwt.RegisteredClaims
}

func GenerateJWT(user model.User) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		IsPremium: user.Profile.IsPremium,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "kopelko-dating-app-backend",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parsing JWT: %w", err)
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func LoadJWTKey() []byte {
	jwtKey = []byte(os.Getenv("JWT_KEY"))
	return jwtKey
}
