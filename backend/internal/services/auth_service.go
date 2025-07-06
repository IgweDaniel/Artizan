package services

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	repoInterfaces "github.com/igwedaniel/artizan/internal/interfaces/repositories"
	"github.com/igwedaniel/artizan/internal/models"
	"github.com/igwedaniel/artizan/pkg/utils"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrExpiredNonce     = errors.New("nonce has expired")
	ErrInvalidJWT       = errors.New("invalid JWT token")
)

// constant prefix for access and refresh token key=> prefix(constant) + secret(dynamic)

const (
	accessTokenPrefix    = "access_token_"
	accessTokenDuration  = 24 * time.Hour
	refreshTokenPrefix   = "refresh_token_"
	refreshTokenDuration = 30 * 24 * time.Hour // 30 days
)

const (
	messageValidity = 5 * time.Minute
	nonceLength     = 32
	messageTemplate = "Welcome to Artizan!\n\nPlease sign this message to verify your wallet ownership.\n\nNonce: %s\nAddress: %s"
)

type AuthService struct {
	userRepo      repoInterfaces.UserRepository
	authNonceRepo repoInterfaces.AuthNonceRepository
	secret        string
}

// NewAuthService creates a new instance of AuthService
func NewAuthService(secret string, userRepo repoInterfaces.UserRepository, authNonceRepo repoInterfaces.AuthNonceRepository) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		authNonceRepo: authNonceRepo,
		secret:        secret,
	}
}

// authenticate using signature and nonce
func (s *AuthService) GetNonceMessage(walletAddress string) (string, error) {
	authNonce, err := s.authNonceRepo.GetByAddress(walletAddress)
	if err != nil {
		if !errors.Is(err, repoInterfaces.ErrRecordNotFound) {
			return "", err
		}
		// Not found, create new nonce
		return s.createAndStoreNonceMessage(walletAddress)
	}

	// Check if the nonce has expired
	if time.Since(authNonce.CreatedAt) > messageValidity {
		return s.createAndStoreNonceMessage(walletAddress)
	}

	return authNonce.Message, nil
}

// Helper to create and store a new nonce message
func (s *AuthService) createAndStoreNonceMessage(walletAddress string) (string, error) {
	nonce, err := utils.Randomize(nonceLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	message := fmt.Sprintf(messageTemplate, nonce, walletAddress)
	authNonce := &models.AuthNonce{
		WalletAddress: walletAddress,
		Message:       message,
	}
	if err := s.authNonceRepo.Create(authNonce); err != nil {
		return "", err
	}
	return message, nil
}

// authenticate using signature and nonce
func (s *AuthService) Authenticate(walletAddress, signature string) (*AuthTokens, *models.User, error) {
	nonceMsg, err := s.authNonceRepo.GetByAddress(walletAddress)
	if err != nil {
		if errors.Is(err, repoInterfaces.ErrRecordNotFound) {
			return nil, nil, ErrInvalidSignature
		}
		return nil, nil, err
	}

	if err := s.verifySignature(walletAddress, signature, nonceMsg.Message); err != nil {
		return nil, nil, ErrInvalidSignature
	}

	// Check if the nonce message is still valid
	if time.Since(nonceMsg.CreatedAt) > messageValidity {
		return nil, nil, ErrExpiredNonce
	}

	user, err := s.userRepo.GetUserByWalletAddress(walletAddress)
	if err != nil {
		if !errors.Is(err, repoInterfaces.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("failed to get user: %w", err)
		}
		// User not found, create a new user
		user = &models.User{
			WalletAddress: walletAddress,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	tokens, err := s.generateTokens(user.WalletAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, user, nil
}

// verifyToken verifies the provided token and returns the wallet address if valid

func (s *AuthService) VerifyToken(token string) (*models.User, error) {
	claims, err := utils.ParseJWT(token, []byte(s.secret+accessTokenPrefix))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w %w", err, ErrInvalidJWT)
	}
	if claims == nil || claims.ID == "" {
		return nil, ErrInvalidJWT
	}

	user, err := s.userRepo.GetUserByWalletAddress(claims.ID)
	if err != nil {
		if errors.Is(err, repoInterfaces.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found for wallet address %s: %w", claims.ID, ErrInvalidJWT)
		}
		return nil, fmt.Errorf("failed to get user by wallet address %s: %w", claims.ID, err)
	}

	return user, nil
}

// refreshToken refreshes the provided token and returns a new token
func (s *AuthService) RefreshToken(token string) (*AuthTokens, *models.User, error) {

	claims, err := utils.ParseJWT(token, []byte(s.secret+refreshTokenPrefix))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse refresh token: %w %w", err, ErrInvalidJWT)
	}
	if claims == nil || claims.ID == "" {
		return nil, nil, ErrInvalidJWT
	}
	user, err := s.userRepo.GetUserByWalletAddress(claims.ID)
	if err != nil {
		if errors.Is(err, repoInterfaces.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user not found for wallet address %s: %w", claims.ID, ErrInvalidJWT)
		}
		return nil, nil, fmt.Errorf("failed to get user by wallet address %s: %w", claims.ID, err)
	}

	tokens, err := s.generateTokens(user.WalletAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, user, nil
}

// generateTokens creates access and refresh tokens for a wallet address.
func (s *AuthService) generateTokens(walletAddress string) (*AuthTokens, error) {
	accessToken, err := utils.IssueJWT(walletAddress, accessTokenDuration, []byte(s.secret+accessTokenPrefix))
	if err != nil {
		return nil, fmt.Errorf("failed to issue access token: %w", err)
	}
	refreshToken, err := utils.IssueJWT(walletAddress, refreshTokenDuration, []byte(s.secret+refreshTokenPrefix))
	if err != nil {
		return nil, fmt.Errorf("failed to issue refresh token: %w", err)
	}
	return &AuthTokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func (s *AuthService) verifySignature(address, signature, message string) error {
	prefix := "\x19Ethereum Signed Message:\n"
	dataLength := strconv.Itoa(len(message))
	formattedMessage := []byte(prefix + dataLength + message)
	msgHash := crypto.Keccak256Hash([]byte(formattedMessage))

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	if len(sig) != 65 {
		return fmt.Errorf("invalid signature length")
	}

	if sig[64] >= 27 {
		sig[64] -= 27
	}

	// Recover public key from signature
	pubKey, err := crypto.SigToPub(msgHash.Bytes(), sig)
	if err != nil {
		return fmt.Errorf("failed to recover public key from signature: %v", err)
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()
	if !strings.EqualFold(recoveredAddress, address) {
		return fmt.Errorf("address does not match the public key derived from the signature")
	}

	return nil
}
