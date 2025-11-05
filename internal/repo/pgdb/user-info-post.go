package pgdb

import (
	"context"
	"fmt"
	"sync"

	"support-bot/internal/entity"
	"support-bot/internal/repo/repoerrors"
	"support-bot/pkg/postgres"
)

type RepoUserInfoPost struct {
	*postgres.Postgres
	mu *sync.RWMutex
}

func NewRepoUserInfoPost(pg *postgres.Postgres) *RepoUserInfoPost {
	return &RepoUserInfoPost{pg, &sync.RWMutex{}}
}

func (r *RepoUserInfoPost) Add(ctx context.Context, entity entity.UserInfoPost) error {
	sql, args, _ := r.Builder.
		Insert("user_info_posts").
		Columns("user_id", "forward_post_id").
		Values(entity.UserId, entity.ForwardPostId).
		ToSql()

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return repoerrors.ErrNotInserted
	}

	return nil
}

// func (r *RepoUserInfoPost) Edit(ctx context.Context, postId int64, fwrdId int64) error {
// 	sql, args, _ := r.Builder.
// 		Update("user_info_posts").
// 		Set("forward_id", fwrdId).
// 		Where("post_id = ?", postId).
// 		ToSql()

// 	result, err := r.Pool.Exec(ctx, sql, args...)
// 	if err != nil {
// 		return err
// 	}

// 	if result.RowsAffected() != 1 {
// 		return repoerrors.ErrNotInserted
// 	}

// 	return nil
// }

func (r *RepoUserInfoPost) GetList(ctx context.Context) (map[int64]int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list = make(map[int64]int64)

	sql, args, _ := r.Builder.
		Select("user_id, forward_post_id").
		From("user_info_posts").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var record entity.UserInfoPost
		err = rows.Scan(&record.UserId, &record.ForwardPostId)
		if err != nil {
			return list, fmt.Errorf("RepoMessageUpdate.GetList - rows.Scan: %v", err)
		}
		list[record.UserId] = record.ForwardPostId
	}

	return list, nil
}
