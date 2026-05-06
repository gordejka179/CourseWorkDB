package usecase

//Возврат книги
func (s *Service) ReturnBook(readerLibraryCard string, inventoryNumber string) error{
    return s.repo.ReturnBook(readerLibraryCard, inventoryNumber)
}
