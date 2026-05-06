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

type Building struct {
    BuildingId int
    Address string
}

type Publication struct {
    ID int
    Title string
    PublicationYear int
    Authors []Author
    ISBNs []string
    BBKs []string
    OtherIndexes []string
}

type Copy struct {
    CopyId int 
    PublicationId int
    BuildingId int
    ReaderId int
    LibrarianId int
    Address string
    Description string
}

type BookingInformation struct {
    CopyId int
    InventoryNumber string
    Title string
    PublicationYear int
    Authors []Author
    Isbns []string
	BBKs []string
    OtherIndexes []string
    Building Building
    ExpiryDate string
}

type IssueInformation struct {
    CopyId int
    ExpiryDate string
    InventoryNumber string
    Title string
    PublicationYear int
    Authors []Author
    Isbns []string
	BBKs []string
    OtherIndexes []string
    Building Building
}