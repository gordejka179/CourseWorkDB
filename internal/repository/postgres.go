package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gordejka179/CourseWorkDB/internal/models"
	_ "github.com/lib/pq"
)


type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func NewPostgresDB(cfg models.DBConfig) (*sql.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBUsername, cfg.DBPassword, cfg.DBName,
    )


    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil {
        db.Close()
        return nil, fmt.Errorf("db unreachable: %w", err)
    }

    log.Println("PostgreSQL connected")
    return db, nil
}