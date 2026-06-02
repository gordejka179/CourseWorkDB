package usecase

import "github.com/gordejka179/CourseWorkDB/internal/models"

//Возврат книги
func (s *Service) ReturnBook(inventoryNumber string) error{
    return s.repo.ReturnBook(inventoryNumber)
}


func (s *Service) GetBuildings()([]models.Building2, error) {
    return s.repo.GetBuildings()
}