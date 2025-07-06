package main

import (
	"fmt"
	"log"

	"github.com/igwedaniel/artizan/internal/adapters/eventbus"
	"github.com/igwedaniel/artizan/internal/adapters/http"
	"github.com/igwedaniel/artizan/internal/adapters/repositories"
	"github.com/igwedaniel/artizan/internal/config"
	"github.com/igwedaniel/artizan/internal/eventhandlers"
	"github.com/igwedaniel/artizan/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DbUrl), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := repositories.NewGormUserRepository(db)
	authNonceRepo := repositories.NewGormAuthNonceRepository(db)
	eventBus := eventbus.New()

	svcs := &http.Services{
		AuthService: services.NewAuthService(cfg.JwtSecret, userRepo, authNonceRepo),
		UserService: services.NewUserService(userRepo),
	}

	// Example: subscribe to a user.created event
	eventBus.Subscribe(eventhandlers.EventUserCreated, eventhandlers.NewHandleUserCreatedEvent(svcs.UserService))

	e := http.NewServer(svcs)

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
