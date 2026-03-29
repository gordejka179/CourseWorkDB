package config

import "github.com/gordejka179/CourseWorkDB/internal/models"


func InitConfig() models.Config {
	dbC := models.DBConfig{
		DBUsername:    "postgres",
		DBPassword:    "qwerty",
		DBName:        "library",
		DBHost:        "localhost",
	}

	sC := models.ServerConfig{
		Port : "8080",
	}
	
	cfg := models.Config{
		Server : sC,
		DB : dbC,

	}

	return cfg
}