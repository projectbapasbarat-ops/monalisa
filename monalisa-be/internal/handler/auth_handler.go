package handler

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"monalisa-be/internal/repository"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(u *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: u}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		NIP string `json:"nip"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.NIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	// =========================
	// DEBUG JWT SECRET (SIGN)
	// =========================
	log.Println("JWT_SECRET (auth):", os.Getenv("JWT_SECRET"))

	auth, err := h.userRepo.GetUserAuthByNIP(req.NIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	claims := jwt.MapClaims{
		"user_id":     auth.UserID,
		"permissions": auth.Permissions,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signed,
	})
}
