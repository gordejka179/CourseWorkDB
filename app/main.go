package main

import (
	"context"
	"log"

	_ "github.com/gordejka179/repository"
)

type Config struct {
	appPort string
}

func main() {
	conf := InitConfig()


	ctx := context.Background()
    pool, err := db.NewPostgresPool(ctx, "postgres://user:pass@localhost:5432/library?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    repo := repository.NewBookRepository(pool)

    reserveUseCase := usecase.NewReserveBook(repo)

    handler := handler.NewLibraryHandler(reserveUseCase)


	if err := srv.Run(conf.appPort, handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}

func InitConfig() Config {
	return Config{
		appPort: "8080",
	}
}