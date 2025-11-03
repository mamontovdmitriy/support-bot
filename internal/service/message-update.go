package service

import (
	"context"
	"support-bot/internal/entity"
	"support-bot/internal/repo"
)

type MessageUpdateService struct {
	repo repo.MessageUpdate
}

func NewMesssageUpdateService(repo repo.MessageUpdate) *MessageUpdateService {
	return &MessageUpdateService{repo: repo}
}

func (s *MessageUpdateService) Create(ctx context.Context, entity entity.MessageUpdate) (int, error) {
	return s.repo.Create(ctx, entity)
}
