package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

type Repository interface {
    CheckIfReaderExists(login string) (bool, error)
    CreateReader(reader *models.Reader) error
    CheckReaderCredentials(email, password string) (bool, error)
    CheckLibrarianCredentials(email, password string) (success bool, err error)

    GetPublicationsByISBN(ISBN string) ([]models.Publication, error)
    GetPublicationsByTitle(Title string) ([]models.Publication, error)
    GetPublicationsByAuthor(Author models.Author) ([]models.Publication, error)
    GetFullCodes(bbkCodes []string) ([]string, error)
    GetAdditionalCodes(fullCodes []string) ([]string, error)
    GetPublicationsByBBK(allCodes []string) ([]models.Publication, error)


    GetCopiesByIDList(ids []int)([]models.Copy, error)

    ReserveCopyByEmail(email string, copyId int)(error)

    GetCurrentBookingsByEmail(email string) ([]models.BookingInformation, error)

    MakeLoan(emailLibrarian string, emailReader string, inventoryNumber string) error
    GetLoanedBooksByReaderEmail(emailReader string) ([]models.IssueInformation, error)

    GetLoanedBooksByReaderLibraryCard(readerLibraryCard string) ([]models.IssueInformation, error)

    GetCurrentBookingsByReaderLibraryCard(readerLibraryCard string) ([]models.BookingInformation, error)

    ReturnBook(readerLibraryCard string, inventoryNumber string) error
}

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}


func (s *Service) CheckIfReaderExists(email string) (bool, error) {
    return s.repo.CheckIfReaderExists(email)
}

func (s *Service) CreateReader(reader *models.Reader) error {
    return s.repo.CreateReader(reader)
}

func (s *Service) CheckReaderCredentials(email, password string) (success bool, err error) {
    return s.repo.CheckReaderCredentials(email, password)
}

func (s *Service) CheckLibrarianCredentials(email, password string) (success bool, err error) {
    return s.repo.CheckLibrarianCredentials(email, password)
}

