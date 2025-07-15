package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ttodoshi/code-typing-auth-service/internal/adapters/handler/http/api"
	"github.com/ttodoshi/code-typing-auth-service/pkg/logging"
	"net/http"
)

type Router struct {
	log logging.Logger
	*api.AuthHandler
}

func NewRouter(log logging.Logger, authHandler *api.AuthHandler) *Router {
	return &Router{
		log:         log,
		AuthHandler: authHandler,
	}
}

func (r *Router) InitRoutes(e *gin.Engine) {
	r.log.Info("initializing error handling middleware")
	e.Use(ErrorHandlerMiddleware())
	e.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.log.Info("initializing routes")

	// swagger
	e.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// healthcheck
	e.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	apiGroup := e.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")

	v1TextsGroup := v1ApiGroup.Group("/auth")
	{
		v1TextsGroup.POST("/registration", r.Register)
		v1TextsGroup.POST("/login", r.Login)
		v1TextsGroup.GET("/refresh", r.Refresh)
		v1TextsGroup.DELETE("/logout", r.Logout)
	}
}
