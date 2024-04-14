package main

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/igostfost/avito_backend_trainee"
	"github.com/igostfost/avito_backend_trainee/pkg/handler"
	"github.com/igostfost/avito_backend_trainee/pkg/repository"
	"github.com/igostfost/avito_backend_trainee/pkg/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
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

	db, err := connectToDB()
	if err != nil {
		logrus.Printf("error init data base: %s", err)
	}
	defer db.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pass"),
		DB:       viper.GetInt("redis.db"),
	})

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		logrus.Fatalf("error connecting to Redis: %s", err)
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

func connectToDB() (*sqlx.DB, error) {
	cfg := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	maxAttempts := 5
	attempt := 1
	var db *sqlx.DB
	var err error
	for attempt <= maxAttempts {
		db, err = repository.NewPostgresDB(cfg)
		if err == nil {
			return db, nil
		}

		time.Sleep(3 * time.Second)
		attempt++
	}

	return nil, errors.New("failed to connect to the database after multiple attempts")
}
