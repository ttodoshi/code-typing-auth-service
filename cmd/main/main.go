package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	_ "speed-typing-auth-service/docs"
	"speed-typing-auth-service/internal/adapters/handler"
	"speed-typing-auth-service/internal/adapters/mq/rabbitmq"
	"speed-typing-auth-service/internal/adapters/repository/mongodb"
	"speed-typing-auth-service/internal/core/servises"
	"speed-typing-auth-service/pkg/broker"
	"speed-typing-auth-service/pkg/discovery"
	"speed-typing-auth-service/pkg/env"
	"speed-typing-auth-service/pkg/logging"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

var (
	authHandler *handler.AuthHandler
	log         logging.Logger
)

func init() {
	env.LoadEnvVariables()
	if os.Getenv("PROFILE") == Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	log = logging.GetLogger()
	discovery.InitServiceDiscovery()
	err := mgm.SetDefaultConfig(nil, "auth", options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		log.Fatal("failed connect to database")
	}
}

func main() {
	refreshTokenRepository := mongodb.NewRefreshTokenRepository()
	userRepository := mongodb.NewUserRepository()

	channel := broker.InitMessageBroker()
	defer broker.Close()

	resultsMigrator := rabbitmq.NewResultsMigrator(channel, log)
	authService := servises.NewAuthService(
		userRepository, refreshTokenRepository,
		resultsMigrator,
		log,
	)
	authHandler = handler.NewAuthHandler(authService, log)

	initRoutes()
}

//	@title						Auth Service API
//	@version					1.0

//	@host						localhost:8090
//	@BasePath					/api/v1

// @externalDocs.description	OpenAPI
func initRoutes() {
	r := gin.Default()

	log.Info("initializing error handling middleware")
	r.Use(handler.ErrorHandlerMiddleware())

	log.Info("initializing handlers")

	// swagger
	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := r.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")

	v1TextsGroup := v1ApiGroup.Group("/auth")
	{
		v1TextsGroup.POST("/registration", authHandler.Register)
		v1TextsGroup.POST("/login", authHandler.Login)
		v1TextsGroup.GET("/refresh", authHandler.Refresh)
		v1TextsGroup.DELETE("/logout", authHandler.Logout)
	}

	log.Infof("starting server on port :%s", os.Getenv("PORT"))

	err := r.Run()
	if err != nil {
		log.Fatal("error while running server")
	}
}
