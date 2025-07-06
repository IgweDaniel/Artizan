package interfaces

import "github.com/igwedaniel/artizan/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetUserByWalletAddress(walletAddress string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	UpdateUserByID(id string, user *models.User) error
	DeleteUserByID(id string) error
}
