package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"


func (s *Service) ReserveCopyByEmail(email string, copyId int)(error){
    return s.repo.ReserveCopyByEmail(email, copyId)
}

func (s *Service) GetCurrentBookingsByEmail(email string) ([]models.BookingInformation, error) {
    return s.repo.GetCurrentBookingsByEmail(email)
}

func (s *Service) GetCurrentBookingsByReaderLibraryCard (readerLibraryCard string) ([]models.BookingInformation, error){
    return s.repo.GetCurrentBookingsByReaderLibraryCard(readerLibraryCard)
}