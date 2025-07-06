package services

import (
	"errors"

	repoInterfaces "github.com/igwedaniel/artizan/internal/interfaces/repositories"
	"github.com/igwedaniel/artizan/internal/models"
)

type UserService struct {
	userRepo repoInterfaces.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repoInterfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// get User by ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repoInterfaces.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// update User By ID
func (s *UserService) UpdateUserByID(id string, user *models.User) error {
	// implementation
	return nil
}

// delete User By ID
func (s *UserService) DeleteUserByID(id string) error {
	// implementation
	return nil
}

// get all Users with options
func (s *UserService) GetAllUsers() ([]*models.User, error) {
	// implementation
	return nil, nil
}

// get all creators
func (s *UserService) GetAllCreators() ([]*models.User, error) {
	// implementation
	return nil, nil
}

// get top creators
func (s *UserService) GetTopCreators(limit int) ([]*models.User, error) {
	// implementation
	return nil, nil
}
