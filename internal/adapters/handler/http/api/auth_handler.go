package api

import (
	"code-typing-auth-service/internal/core/ports"
	"code-typing-auth-service/internal/core/ports/dto"
	"code-typing-auth-service/internal/core/ports/errors"
	"code-typing-auth-service/pkg/jwt"
	"code-typing-auth-service/pkg/logging"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	cookieHost = os.Getenv("COOKIE_HOST")
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

// Register godoc
//
//	@Summary		Register new user
//	@Description	Register new user
//	@Tags			auth
//	@Accept			json
//	@Produce		plain
//	@Param			request	body		dto.RegisterRequestDto	true	"Register request"
//	@Success		201		{object}	string
//	@Header			200		{string}	Set-Cookie	"refreshToken"
//	@Router			/auth/registration [post]
func (h *AuthHandler) Register(c *gin.Context) {
	h.log.Debug("received register request")

	sessionCookie, err := c.Cookie("SESSION")
	var registerRequestDto dto.RegisterRequestDto
	if err = c.ShouldBindJSON(&registerRequestDto); err != nil {
		err = c.Error(&errors.BodyMappingError{
			Message: "error in request body",
		})
		h.log.Warn("error in request body")
		return
	}

	access, refresh, err := h.svc.Register(registerRequestDto, sessionCookie)
	if err != nil {
		err = c.Error(err)
		return
	}

	c.SetCookie("refreshToken", refresh, jwt.RefreshTokenExp, "/", cookieHost, false, true)
	c.Data(201, "text/html; charset=utf-8", []byte(access))
}

// Login godoc
//
//	@Summary		Login
//	@Description	Login
//	@Tags			auth
//	@Accept			json
//	@Produce		plain
//	@Param			request	body		dto.LoginRequestDto	true	"Login request"
//	@Success		200		{object}	string
//	@Header			200		{string}	Set-Cookie	"refreshToken"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	h.log.Debug("received login request")

	sessionCookie, err := c.Cookie("SESSION")
	var loginRequestDto dto.LoginRequestDto
	if err = c.ShouldBindJSON(&loginRequestDto); err != nil {
		err = c.Error(&errors.BodyMappingError{
			Message: "error in request body",
		})
		h.log.Warn("error in request body")
		return
	}

	access, refresh, err := h.svc.Login(loginRequestDto, sessionCookie)
	if err != nil {
		err = c.Error(err)
		return
	}

	c.SetCookie("refreshToken", refresh, jwt.RefreshTokenExp, "/", cookieHost, false, true)
	c.Data(200, "text/html; charset=utf-8", []byte(access))
}

// Refresh godoc
//
//	@Summary		Refresh
//	@Description	Refresh
//	@Tags			auth
//	@Accept			json
//	@Produce		plain
//	@Param			Cookie	header		string	true	"refreshToken"	default(refreshToken=)
//	@Success		200		{object}	string
//	@Header			200		{string}	Set-Cookie	"refreshToken"
//	@Router			/auth/refresh [get]
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

	access, refresh, err := h.svc.Refresh(refreshTokenCookie)
	if err != nil {
		err = c.Error(err)
		return
	}

	c.SetCookie("refreshToken", refresh, jwt.RefreshTokenExp, "/", cookieHost, false, true)
	c.Data(200, "text/html; charset=utf-8", []byte(access))
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Logout
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Cookie	header	string	true	"refreshToken"	default(refreshToken=)
//	@Success		204
//	@Header			204	{string}	Set-Cookie	"refreshToken"
//	@Router			/auth/logout [delete]
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

	h.svc.Logout(refreshTokenCookie)

	c.SetCookie("refreshToken", "", -1, "/", cookieHost, false, true)
	c.Status(204)
}
