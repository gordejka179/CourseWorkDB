package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

func (s *Service) GetOverdueCopies() ([]models.OverdueCopy, error) {
    return s.repo.GetOverdueCopies()
}

func (s *Service) GetOverallStats() (*models.OverallStats, error) {
    return s.repo.GetOverallStats()
}