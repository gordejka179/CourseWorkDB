package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

//сделать выдачу
func (s *Service) MakeLoan(readerLibraryCard string, emailLibrarian string, inventoryNumber string) error{
    return s.repo.MakeLoan(readerLibraryCard, emailLibrarian, inventoryNumber)
}

//получить информацию о выданных читателю книгах
func (s *Service) GetLoanedBooksByReaderEmail(emailReader string) ([]models.IssueInformation, error) {
	return s.repo.GetLoanedBooksByReaderEmail(emailReader)
}

func (s *Service) GetLoanedBooksByReaderLibraryCard(readerLibraryCard string) ([]models.IssueInformation, error) {
	return s.repo.GetLoanedBooksByReaderLibraryCard(readerLibraryCard)
}




