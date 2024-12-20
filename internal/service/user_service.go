package service

import (
	"myapp/internal/domain"
	"myapp/internal/repository"
)

type UserService interface {
	GetUserByID(id int64) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(id int64) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) GetUserByID(id int64) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) CreateUser(user *domain.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) UpdateUser(user *domain.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id int64) error {
	return s.userRepo.Delete(id)
}
