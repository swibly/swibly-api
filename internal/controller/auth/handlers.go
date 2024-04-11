package auth

import (
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(ctx *gin.Context) {
	var body struct {
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	userID, err := register(body.Fullname, body.Username, body.Email, body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	token, err := utils.GenerateJWT(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(ctx *gin.Context) {
	var body struct {
		Username string `json:"username" default:""`
		Email    string `json:"email" default:""`
		Password string `json:"password" default:""`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad body formatting"})
		return
	}

	userID, err := login(body.Username, body.Email, body.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	token, err := utils.GenerateJWT(userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
