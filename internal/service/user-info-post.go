package service

import (
	"context"
	"errors"
	"support-bot/internal/entity"
	"support-bot/internal/repo"
	"sync"
)

type UserInfoPostService struct {
	repo repo.UserInfoPost
	mu   sync.RWMutex
	list map[int64]int64
}

func NewUserInfoPostService(repo repo.UserInfoPost) *UserInfoPostService {
	return &UserInfoPostService{repo: repo}
}

func (s *UserInfoPostService) SaveUserInfoPost(userId int64, forwardPostId int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.repo.Add(ctx, entity.UserInfoPost{UserId: userId, ForwardPostId: forwardPostId})
	if err != nil {
		return err
	}

	s.list[userId] = forwardPostId

	return nil
}

func (s *UserInfoPostService) GetForwardId(userId int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(s.list) == 0 {
		list, _ := s.repo.GetList(ctx)
		s.list = list
	}

	forwardPostId, exists := s.list[userId]
	if !exists {
		return 0, errors.New("ForwardPostID does not exist")
	}

	return forwardPostId, nil
}

func (s *UserInfoPostService) GetUserId(forwardId int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(s.list) == 0 {
		list, _ := s.repo.GetList(ctx)
		s.list = list
	}

	// поиск использует меньше итераций, чем разворачивание мапы
	for k, v := range s.list {
		if v == forwardId {
			return k, nil
		}
	}

	return 0, errors.New("UserId does not exist")
}
