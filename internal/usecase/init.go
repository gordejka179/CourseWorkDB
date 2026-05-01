package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

type Repository interface {
    CheckIfReaderExists(login string) (bool, error)
    CreateReader(reader *models.Reader) error
    CheckReaderCredentials(email, password string) (bool, error)
    CheckLibrarianCredentials(email, password string) (success bool, err error)
    GetPublicationsByISBN(ISBN string) ([]models.Publication, error)

    GetCopiesByIDList(ids []int)([]models.Copy, error)

    ReserveCopyByEmail(email string, copyId int)(error)

    
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
