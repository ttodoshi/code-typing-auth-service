package main

import (
	_ "code-typing-auth-service/docs"
	"code-typing-auth-service/internal/adapters/handler/http"
	"code-typing-auth-service/internal/adapters/handler/http/api"
	"code-typing-auth-service/internal/adapters/mq/rabbitmq"
	"code-typing-auth-service/internal/adapters/repository/mongodb"
	"code-typing-auth-service/internal/core/domain"
	"code-typing-auth-service/internal/core/servises"
	"code-typing-auth-service/pkg/broker"
	"code-typing-auth-service/pkg/discovery"
	"code-typing-auth-service/pkg/env"
	"code-typing-auth-service/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strconv"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

func init() {
	env.LoadEnvVariables()
	if os.Getenv("PROFILE") == Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	discovery.InitServiceDiscovery()
}

//	@title		Auth Service API
//	@version	1.0

// @host		localhost:8090
// @BasePath	/api/v1
func main() {
	log := logging.GetLogger()

	initDatabase(log)

	channel := broker.InitMessageBroker()
	defer broker.Close()

	r := gin.Default()
	router := initRouter(log, channel)
	router.InitRoutes(r)

	log.Fatalf("error while running server due to: %s", r.Run())
}

func initDatabase(log logging.Logger) {
	err := mgm.SetDefaultConfig(nil, "auth", options.Client().ApplyURI(os.Getenv("DB_URL")))

	if err != nil {
		log.Fatal("failed connect to database")
	}

	collection := mgm.CollectionByName(
		mgm.CollName(&domain.RefreshToken{}),
	)
	refreshTokenExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		log.Fatal("failed to parse refresh token expiration")
	}
	refreshTokenExpirationIndex := mongo.IndexModel{
		Keys:    bson.D{{"updated_at", 1}},
		Options: options.Index().SetExpireAfterSeconds(int32(refreshTokenExp)),
	}
	_, err = collection.Indexes().CreateOne(mgm.Ctx(), refreshTokenExpirationIndex)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func initRouter(log logging.Logger, channel *amqp.Channel) *http.Router {
	refreshTokenRepository := mongodb.NewRefreshTokenRepository()
	userRepository := mongodb.NewUserRepository()

	resultsMigrator := rabbitmq.NewResultsMigrator(channel, log)
	authService := servises.NewAuthService(
		userRepository, refreshTokenRepository,
		resultsMigrator,
		log,
	)
	return http.NewRouter(
		log,
		api.NewAuthHandler(
			authService, log,
		),
	)
}
