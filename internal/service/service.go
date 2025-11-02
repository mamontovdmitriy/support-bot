package service

import (
	"support-bot/internal/repo"

	"github.com/sirupsen/logrus"
)

type (
	ServicesDependencies struct {
		Log   *logrus.Logger
		Repos *repo.Repositories
		// other components
	}

	Services struct {
		// Auth        Auth
	}
)

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		// Auth: NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
	}
}
