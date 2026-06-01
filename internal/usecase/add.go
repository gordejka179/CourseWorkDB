package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

func (s *Service) SearchAuthors(lastName string, firstName string, patronymic string, birthYear string) ([]models.AuthorForAdd, error) {
    return s.repo.SearchAuthors(lastName, firstName, patronymic, birthYear)
}

func (s *Service) CreateAuthor(lastName, firstName, patronymic, birthYear string) error {
    return s.repo.CreateAuthor(lastName, firstName, patronymic, birthYear)
}

func (s *Service) CreatePublication(title string, publicationYear int, authorIds []int, isbns []string, otherIsbns []string, bbks []string, otherIndexes []string) error {
    return s.repo.CreatePublication(title, publicationYear, authorIds, isbns, otherIsbns, bbks, otherIndexes)
}
