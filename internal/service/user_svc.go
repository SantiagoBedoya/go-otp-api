package service

import "github.com/SantiagoBedoya/otp-api/internal/model"

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByID(userID string) (*model.User, error)
	Save(u *model.User) error
	SaveSecret(userID, secret string) error
	SetValidSecret(userID string) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUsers() ([]model.User, error) {
	return s.repo.GetAll()
}
