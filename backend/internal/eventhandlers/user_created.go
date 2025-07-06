package eventhandlers

import (
	interfaces "github.com/igwedaniel/artizan/internal/interfaces/eventbus"
	"github.com/igwedaniel/artizan/internal/services"
)

func NewHandleUserCreatedEvent(UserService *services.UserService) func(interfaces.Event) {
	return func(e interfaces.Event) {
		// handler logic using userRepo
	}
}
