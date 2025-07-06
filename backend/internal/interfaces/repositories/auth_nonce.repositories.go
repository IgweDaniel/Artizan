package interfaces

import "github.com/igwedaniel/artizan/internal/models"

type AuthNonceRepository interface {
	Create(authNonce *models.AuthNonce) error
	GetByAddress(walletAddress string) (*models.AuthNonce, error)
}
