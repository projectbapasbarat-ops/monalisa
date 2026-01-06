package handler

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"monalisa-be/internal/auth"
	"monalisa-be/internal/repository"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.RefreshTokenRepository
}

func NewAuthHandler(
	userRepo *repository.UserRepository,
	tokenRepo *repository.RefreshTokenRepository,
) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

/*
LOGIN
- input: NIP (password menyusul di step berikutnya)
- output: access_token + refresh_token
*/
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		NIP string `json:"nip"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.NIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	authData, err := h.userRepo.GetUserAuthByNIP(req.NIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(
		authData.UserID,
		authData.Permissions,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate access token"})
		return
	}

	refreshDays, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL_DAYS"))
	refreshToken, expiresAt, err := auth.GenerateRefreshToken(
		authData.UserID,
		refreshDays,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate refresh token"})
		return
	}

	// simpan refresh token ke DB
	_ = h.tokenRepo.Save(authData.UserID, refreshToken, expiresAt)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

/*
REFRESH ACCESS TOKEN
*/
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	userID, err := h.tokenRepo.FindValid(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid refresh token"})
		return
	}

	permissions, err := h.userRepo.GetPermissionsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to load permissions"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(userID, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

/*
LOGOUT
*/
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err == nil {
		_ = h.tokenRepo.Revoke(req.RefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
