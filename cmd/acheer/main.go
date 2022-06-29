package main

import (
	"context"
	"github.com/miprokop/fication/configs"
	server "github.com/miprokop/fication/internal/http-server"
	"github.com/miprokop/fication/internal/http-server/handlers"
	"github.com/miprokop/fication/internal/persistence/postgres"
	"github.com/miprokop/fication/internal/services"
	"github.com/spf13/viper"
	"log"
	"sync"
)

func main() {
	if err := configs.Init("./configs"); err != nil {
		log.Fatalf("error in config init: %s", err)
	}

	db, err := postgres.New(context.Background(), &sync.WaitGroup{}, &configs.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	})

	if err != nil {
		log.Fatalf("error in database init: %s", err)
	}
	rep, err := postgres.NewRepository(db)
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := services.NewService(rep)
	h := handlers.NewHandler(s)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), h.InitRoutes()); err != nil {
		log.Fatalf("err %v in running server", err)
	}
}
