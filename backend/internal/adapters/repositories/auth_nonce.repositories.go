package repositories

import (
	"errors"

	repoInterfaces "github.com/igwedaniel/artizan/internal/interfaces/repositories"
	"github.com/igwedaniel/artizan/internal/models"
	"gorm.io/gorm"
)

type gormAuthNonceRepository struct {
	db *gorm.DB
}

func NewGormAuthNonceRepository(db *gorm.DB) repoInterfaces.AuthNonceRepository {
	return &gormAuthNonceRepository{db: db}
}

func (r *gormAuthNonceRepository) Create(authNonce *models.AuthNonce) error {
	if err := r.db.Create(authNonce).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormAuthNonceRepository) GetByAddress(address string) (*models.AuthNonce, error) {
	var authNonce models.AuthNonce
	if err := r.db.Where("wallet_address = ?", address).First(&authNonce).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrRecordNotFound
		}
		return nil, err
	}
	return &authNonce, nil
}
