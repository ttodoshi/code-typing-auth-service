package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"speed-typing-auth-service/internal/adapters/dto"
	"speed-typing-auth-service/internal/core/errors"
	"speed-typing-auth-service/internal/core/ports"
	"speed-typing-auth-service/pkg/logging"
	"time"
)

type AuthHandler struct {
	svc ports.AuthService
	log logging.Logger
}

func NewAuthHandler(svc ports.AuthService, log logging.Logger) *AuthHandler {
	return &AuthHandler{
		svc: svc,
		log: log,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	h.log.Debug("received register request")

	var registerRequestDto dto.RegisterRequestDto
	if err := c.ShouldBindJSON(&registerRequestDto); err != nil {
		err = c.Error(&errors.BodyMappingError{
			Message: "error in request body",
		})
		h.log.Warn("error in request body")
		return
	}

	access, refresh, err := h.svc.Register(registerRequestDto)
	if err != nil {
		err = c.Error(err)
		return
	}
	c.JSON(201, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.log.Debug("received login request")

	var loginRequestDto dto.LoginRequestDto
	if err := c.ShouldBindJSON(&loginRequestDto); err != nil {
		err = c.Error(&errors.BodyMappingError{
			Message: "error in request body",
		})
		h.log.Warn("error in request body")
		return
	}

	access, refresh, err := h.svc.Login(loginRequestDto)
	if err != nil {
		err = c.Error(err)
		return
	}

	refreshTokenExp, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXP_SECONDS"))
	if err != nil {
		err = c.Error(&errors.TokenGenerationError{
			Message: fmt.Sprintf("refresh token generation error due to: %s", err.Error()),
		})
		return
	}

	c.SetCookie("refreshToken", refresh, int(refreshTokenExp.Seconds()), "/", os.Getenv("COOKIE_HOST"), false, true)
	c.JSON(200, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	h.log.Debug("received refresh request")

	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil || refreshTokenCookie == "" {
		err = c.Error(&errors.CookieGettingError{
			Message: "error while getting refresh token cookie",
		})
		h.log.Warn("error while getting refreshTokenCookie")
		return
	}

	access, refresh, err := h.svc.Refresh(dto.RefreshRequestDto{
		RefreshToken: refreshTokenCookie,
	})
	if err != nil {
		err = c.Error(err)
		return
	}

	refreshTokenExp, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXP_SECONDS"))
	if err != nil {
		err = c.Error(&errors.TokenGenerationError{
			Message: fmt.Sprintf("refresh token generation error due to: %s", err.Error()),
		})
		return
	}

	c.SetCookie("refreshToken", refresh, int(refreshTokenExp.Seconds()), "/", os.Getenv("COOKIE_HOST"), false, true)
	c.JSON(200, gin.H{
		"access":  access,
		"refresh": refresh,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.log.Debug("received logout request")

	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil || refreshTokenCookie == "" {
		err = c.Error(&errors.CookieGettingError{
			Message: "error while getting refresh token cookie",
		})
		h.log.Warn("error while getting refreshTokenCookie")
		return
	}

	h.svc.Logout(dto.LogoutRequestDto{
		RefreshToken: refreshTokenCookie,
	})

	c.SetCookie("refreshToken", "", 0, "/", os.Getenv("COOKIE_HOST"), false, true)
	c.Status(204)
}
