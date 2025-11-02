package repo

import (
	"support-bot/pkg/postgres"
)

type (
	Repositories struct {
		// ...
	}
)

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		// ...
	}
}
