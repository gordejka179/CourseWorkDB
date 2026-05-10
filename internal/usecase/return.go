package usecase

//Возврат книги
func (s *Service) ReturnBook(inventoryNumber string) error{
    return s.repo.ReturnBook(inventoryNumber)
}
