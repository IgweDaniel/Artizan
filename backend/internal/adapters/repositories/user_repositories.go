package repositories

import (
	"errors"

	repoInterfaces "github.com/igwedaniel/artizan/internal/interfaces/repositories"
	"github.com/igwedaniel/artizan/internal/models"
	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB // Uncomment and use if you have a gorm.DB instance
}

func NewGormUserRepository(db *gorm.DB) repoInterfaces.UserRepository {
	return &gormUserRepository{db: db}
}

// Create creates a new user in the database
func (r *gormUserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// get user by id
func (r *gormUserRepository) GetByID(id string) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrRecordNotFound // Return a specific error if record not found
		}
		return nil, err
	}
	return user, nil
}

// get user by wallet address, case insensitive search
func (r *gormUserRepository) GetUserByWalletAddress(walletAddress string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("LOWER(wallet_address) = ?", walletAddress).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrRecordNotFound // Return a specific error if record not found
		}
		return nil, err
	}
	return &user, nil
}

// update user by id
func (r *gormUserRepository) UpdateUserByID(id string, user *models.User) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrRecordNotFound // Return a specific error if record not found
		}
		return err
	}
	return nil
}

// delete user by id soft delete
func (r *gormUserRepository) DeleteUserByID(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrRecordNotFound // Return a specific error if record not found
		}
		return err
	}
	return nil
}
