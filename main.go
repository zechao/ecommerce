package main

import (
	"log"
	"os"

	"github.com/zechao158/ecomm/cmd/api"
	"github.com/zechao158/ecomm/config"
	"github.com/zechao158/ecomm/storage"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	log.SetFlags()
	log.Println("App running in environment:", config.ENVs.APPEnv)
	db, err := storage.NewPostgreStorage(storage.Config{
		DBUser:     config.ENVs.DBUser,
		DBHost:     config.ENVs.DBHost,
		DBName:     config.ENVs.DBName,
		DBPassword: config.ENVs.DBPassword,
		DBPort:     config.ENVs.DBPort,
		DBSSLMode:  config.ENVs.DBSSLMode,
	})

	initStorage(db)
	server := api.NewAPIServer(config.ENVs.HTTPHost+":"+config.ENVs.HTTPPort, nil)
	err = server.Run()
	if err != nil {
		log.Panicf("error initializing server %v", err)
	}
}

func initStorage(db *gorm.DB) {
	conn, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB connected")
}
