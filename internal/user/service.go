package user

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}
