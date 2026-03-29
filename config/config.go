package config

type Config struct {
    Server ServerConfig
    DB     DBConfig
}

type ServerConfig struct {
    Port string
}

type DBConfig struct {
	dbUsername:    "postgres",
	dbPassword:    "qwerty",
	dbName:        "library",
	dbHost:        "localhost",
}


func InitConfig() Config {
	dbC := DBConfig{
		dbUsername:    "postgres",
		dbPassword:    "qwerty",
		dbName:        "library",
		dbHost:        "localhost",
	}

	sC := ServerConfig{
		app := "8080",
	}
	
	cfg := Config{
		Server := sC,
		DB := dbc,

	}

	return cfg
}