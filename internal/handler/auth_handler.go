package handler

import (
	"auth-service/internal/adapter/blizzard"
	"auth-service/internal/usecase"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	blizzAd blizzard.BlizzardRepository
	uc      usecase.AuthUsecase
	log     *logrus.Logger
}

func NewAuthHandler(
	blizzAd blizzard.BlizzardRepository,
	uc usecase.AuthUsecase,
	log *logrus.Logger,
) *AuthHandler {
	return &AuthHandler{blizzAd: blizzAd, uc: uc, log: log}
}

func (h *AuthHandler) Login(c *gin.Context) {
	stateByte := make([]byte, 7)
	if _, err := rand.Read(stateByte); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate state"})
		return
	}

	state := base64.RawURLEncoding.EncodeToString(stateByte)

	url := h.blizzAd.GetAuthURL(state)
	h.log.WithFields(logrus.Fields{
		"state": state,
	}).Info("Confirm state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Query("code")

	jwtAccess, jwtRefresh, err := h.uc.HandleCallback(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed exchange token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": jwtAccess, "refresh": jwtRefresh})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	jwtToken := c.GetHeader("Authorization")
	refreshStr := strings.TrimPrefix(jwtToken, "Bearer ")
	access, refresh, err := h.uc.RefreshTokens(ctx, refreshStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": access, "refresh": refresh})
}

func (h *AuthHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	blizzID := c.GetString("blizzard_id")
	jwtToken := c.GetHeader("Authorization")
	accessStr := strings.TrimPrefix(jwtToken, "Bearer ")

	user, err := h.uc.GetBlizzardUser(ctx, accessStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": blizzID, "battletag": user.BattleTag})
}
