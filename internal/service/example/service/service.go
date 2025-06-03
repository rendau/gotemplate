package service

import (
	"context"
	"fmt"
)

type Service struct {
	repo RepoI
}

func New(repo RepoI) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Send(ctx context.Context, phone string, msg string, sync bool) error {
	err := s.repo.Send(ctx, phone, msg, sync)
	if err != nil {
		return fmt.Errorf("repo.Send: %w", err)
	}

	return nil
}
