package handler

import (
	"auth-service/pkg/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(
	router *gin.Engine,
	h *AuthHandler,
	cfg *config.Config,
	log *logrus.Logger,
) {
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	router.GET("/auth/login", h.Login)
	router.GET("/callback", h.Callback)

	authProtected := router.Group("/auth")
	authProtected.Use(AuthMiddleware(cfg.JWT.Secret))
	authProtected.POST("/refresh", h.Refresh)
	authProtected.GET("/user", h.GetUser)
}
