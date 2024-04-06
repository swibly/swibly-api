package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/loader"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireAuth(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Parse the tokenString from cookie
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Method of signing in should be HMAC-SHA
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the jwt secret if everything is ok
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		log.Println(err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	// Abort if the token is incorrect
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check if the current time is greater than the expiration date both as float64
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Get the user by the subject "sub" (user.ID)
	var user model.User
	if err := loader.DB.First(&user, claims["sub"]).Error; err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Set a "user" key in the context for Next
	ctx.Set("user", user)

	// Run the next handler
	ctx.Next()
}
