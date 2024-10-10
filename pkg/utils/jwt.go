package utils

import (
	"fmt"
	"time"

	"github.com/swibly/swibly-api/config"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(id uint) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   fmt.Sprint(id),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // Token expires after 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Security.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetClaimsJWT(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Security.JWTSecret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("malformed token: %v", err)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, fmt.Errorf("token is expired or not active yet: %v", err)
			} else {
				return nil, fmt.Errorf("couldn't handle this token: %v", err)
			}
		} else {
			return nil, fmt.Errorf("couldn't handle this token: %v", err)
		}
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token is invalid")
}
