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
}

type Author struct {
    FirstName  string
    LastName   string
    Patronymic string
}

type Publication struct {
    ID int
    Title string
    PublicationYear int
    Authors []Author
    ISBN string
    BBKs []string
    OtherIndexes []string
}

type Copy struct {
    CopyId int 
    InventoryNumber string 
    PublicationId int
    BuildingId int
    ReaderId int
    LibrarianId int
    Address string
    Description string
}