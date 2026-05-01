package usecase


func (s *Service) ReserveCopyByEmail(email string, copyId int)(error){
    return s.repo.ReserveCopyByEmail(email, copyId)
}
