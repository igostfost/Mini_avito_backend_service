package main

import (
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

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error from init configs: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error from load environment variables")
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

	repos := repository.NewRepository(db)
	utilities := utils.NewUtils(repos)
	handlers := handler.NewHandler(utilities)
	serv := new(avito_backend_trainee.Server)

	if err := serv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error running http server: %s", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
