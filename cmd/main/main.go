package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/kamva/mgm/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	_ "speed-typing-auth-service/docs"
	"speed-typing-auth-service/internal/adapters/handler"
	"speed-typing-auth-service/internal/adapters/repository/mongodb"
	"speed-typing-auth-service/internal/core/ports"
	"speed-typing-auth-service/internal/core/servises"
	"speed-typing-auth-service/pkg/env"
	"speed-typing-auth-service/pkg/logging"
	"strconv"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

var (
	authService ports.AuthService
	log         logging.Logger
)

func init() {
	env.LoadEnvVariables()
	if os.Getenv("PROFILE") == Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	log = logging.GetLogger()
	err := mgm.SetDefaultConfig(nil, "auth", options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		log.Fatal("failed connect to database")
	}
}

func main() {
	refreshTokenRepository := mongodb.NewRefreshTokenRepository()
	userRepository := mongodb.NewUserRepository()
	authService = servises.NewAuthService(userRepository, refreshTokenRepository, log)

	initConsul()
	initRoutes()
}

//	@title						Auth Service API
//	@version					1.0
//	@host						localhost:8090
//	@BasePath					/api/v1
//	@externalDocs.description	OpenAPI
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
		authHandler := handler.NewAuthHandler(authService, log)
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

func initConsul() {
	log.Info("initializing consul client")

	consulClient, err := api.NewClient(
		&api.Config{
			Address: os.Getenv("CONSUL_HOST"),
		},
	)
	if err != nil {
		log.Fatal("error creating consul client")
	}

	log.Info("register service in consul")
	agent := consulClient.Agent()
	parsedPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("port parse error")
	}

	service := &api.AgentServiceRegistration{
		Name:    os.Getenv("CONSUL_SERVICE_NAME"),
		Port:    parsedPort,
		Address: os.Getenv("CONSUL_SERVICE_ADDRESS"),
	}
	err = agent.ServiceRegister(service)
	if err != nil {
		log.Fatalf("error while service registration due to error '%s'", err)
	}
}
