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

type UserInfoPost interface {
	Add(ctx context.Context, entity entity.UserInfoPost) error
	// Edit(ctx context.Context, postId int64, fwrdId int64) error
	GetList(ctx context.Context) (map[int64]int64, error)
}

type (
	Repositories struct {
		MessageUpdate
		UserInfoPost
		// ...
	}
)

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		MessageUpdate: pgdb.NewRepoMessageUpdate(pg),
		UserInfoPost:  pgdb.NewRepoUserInfoPost(pg),
		// ...
	}
}
