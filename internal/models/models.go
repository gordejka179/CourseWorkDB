package models

type Config struct {
    Server ServerConfig
    DB     DBConfig
}

type ServerConfig struct {
    Port string
}

type DBConfig struct {
	DBUsername string
	DBPassword string
	DBName string
	DBHost string
}