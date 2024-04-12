package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/igostfost/avito_backend_trainee"
	"github.com/igostfost/avito_backend_trainee/pkg/handler"
	"github.com/igostfost/avito_backend_trainee/pkg/repository"
	"github.com/igostfost/avito_backend_trainee/pkg/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// @title Banners Show App API
// @version 1.0
// @description API Server for BannersShow Application for avito backend trainee

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error from init configs: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error from load environment variables")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	// Проверка соединения с Redis
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("error connecting to Redis: %s", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.SSLMode"),
	})

	if err != nil {
		logrus.Printf("error init data base: %s", err)
	}

	repos := repository.NewRepository(db, redisClient)
	utilities := utils.NewUtils(repos, redisClient)
	handlers := handler.NewHandler(utilities)
	serv := new(avito_backend_trainee.Server)

	if err := serv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error running http server: %s", err.Error())
	}

	logrus.Printf("Server started on port %s", viper.GetString("port"))

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
