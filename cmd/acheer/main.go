package main

import (
	"github.com/miprokop/fication/configs"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := configs.Init("./configs"); err != nil {
		log.Fatalf("error in config init: %s", err)
	}

	db, err := repo.NewPostgresRepository(configs.Config{
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
	rep := repo.NewRepository(db)
	services := service.NewService(rep)
	handlers := router.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitHandlers()); err != nil {
		log.Fatalf("err %v in running server", err)
	}
}
