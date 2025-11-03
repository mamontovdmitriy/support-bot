package repo

import (
	"context"

	"support-bot/internal/entity"
	"support-bot/internal/repo/pgdb"
	"support-bot/pkg/postgres"
)

type MessageUpdate interface {
	Create(ctx context.Context, entity entity.MessageUpdate) (int, error)
	GetById(ctx context.Context, id int) (entity.MessageUpdate, error)
	GetList(ctx context.Context) ([]entity.MessageUpdate, error)
}

type (
	Repositories struct {
		MessageUpdate
		// ...
	}
)

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		MessageUpdate: pgdb.NewRepoMessageUpdate(pg),
		// ...
	}
}
