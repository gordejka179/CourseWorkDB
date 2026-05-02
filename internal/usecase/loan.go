package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

//сделать выдачу
func (s *Service) MakeLoan(readerLibraryCard string, emailLibrarian string, copyId int) error{
    return s.repo.MakeLoan(readerLibraryCard, emailLibrarian, copyId)
}

//получить информацию о выданных читателю книгах
func (s *Service) GetLoanedBooks(emailReader string) ([]models.IssueInformation, error) {
	return s.repo.GetLoanedBooks(emailReader)
}

