package main

import (
	"embed"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/zechao158/ecomm/config"
	"github.com/zechao158/ecomm/storage"
)

//go:embed *.sql
var embedMigrations embed.FS

func main() {
	db, err := storage.NewPostgreStorage(storage.Config{
		DBUser:     config.ENVs.DBUser,
		DBHost:     config.ENVs.DBHost,
		DBName:     config.ENVs.DBName,
		DBPassword: config.ENVs.DBPassword,
		DBPort:     config.ENVs.DBPort,
		DBSSLMode:  config.ENVs.DBSSLMode,
	})
	if err != nil {
		log.Panic(err)
	}

	sqldb, err := db.DB()
	if err != nil {

	}

	log.Println("runnig migration")
	goose.SetBaseFS(embedMigrations)

	dir := ""
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	switch dir {
	case "up":
		if err := goose.Up(sqldb, "."); err != nil {
			log.Panic(err)
		}
	case "down":
		if err := goose.Down(sqldb, "."); err != nil {
			log.Panic(err)
		}
	default:
		log.Fatal("Please, use up or down")
	}

}
