package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

func (s *Service) GetCurrentBookingsByEmail(email string) ([]models.BookingInformation, error) {
    return s.repo.GetCurrentBookingsByEmail(email)
}