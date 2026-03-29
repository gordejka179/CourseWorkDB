package main

import (
	"log"

	server "github.com/gordejka179/CourseWorkDB"
	"github.com/gordejka179/CourseWorkDB/config"
	"github.com/gordejka179/CourseWorkDB/internal/handler"
	"github.com/gordejka179/CourseWorkDB/internal/repository"
	"github.com/gordejka179/CourseWorkDB/internal/usecase"
)


func main() {
	conf := config.InitConfig()


	db, err := repository.NewPostgresDB(conf.DB)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
    defer db.Close()

    repo := repository.NewRepository(db)

    service := usecase.NewService(repo)

    handler := handler.NewUserHandler(service)

	server := server.NewServer()
	if err := server.Run(conf.Server.Port, handler.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}

