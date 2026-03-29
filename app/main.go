package main

import (
	"log"

	_ "github.com/gordejka179/CourseWorkDB/config"
	_ "github.com/gordejka179/CourseWorkDB/internal/handler"
	_ "github.com/gordejka179/CourseWorkDB/internal/usecase"
	"github.com/gordejka179/CourseWorkDB/repository"
	_ "github.com/gordejka179/CourseWorkDB/server"
)


func main() {
	conf := config.InitConfig()

	err := repository.CreateTables(conf.dbUsername, conf.dbPassword, conf.dbHost, conf.dbName)

	if err != nil {
		log.Fatalf("Error occured while creating table: %s%s", err.Error(), conf.appPort, conf.dbUsername, conf.dbPassword, conf.dbName, conf.SessionSecret, conf.dbHost)
	}

	db, err := repository.NewPostgresDB(conf.DB)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
    defer db.Close()

    repo := repository.NewRepository(db)

    service := usecase.NewService(repo)

    handler := handler.NewUserHandler(service)

	if err := server.Run(conf.appPort, handler.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}

