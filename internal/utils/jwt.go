package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID uint) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT secret not set")
	}

	// Generate JWT token
	// "sub" is the subject of the token (ID)
	// "exp" is the expiration date
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Will add 7 days
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))

	// Only related to internal error, not user input
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method used is HMAC, to match the server's signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Use the server's JWT secret for signing
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Error handling for token parsing issues
	if err != nil {
		return nil, err
	}

	// Type assert the token claims to MapClaims for use
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	// If the function reaches this point, it means the type assertion failed.
	// Return nil for the claims and an error indicating the failure.
	return nil, fmt.Errorf("failed to assert token claims")
}
