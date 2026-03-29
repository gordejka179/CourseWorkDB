package usecase

type UserRepository interface {
}

type Service struct {
    repo UserRepository
}

func NewService(repo UserRepository) *Service {
    return &Service{repo: repo}
}