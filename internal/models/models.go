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

type Reader struct{
    ReaderId int
    Email string
    LibraryCard string
    PassportSeries string
    PassportNumber string
    FirstName string
    LastName string
    Patronymic string
    PasswordHash string
};