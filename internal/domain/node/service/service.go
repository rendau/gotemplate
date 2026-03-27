package service

import (
	"context"
	"fmt"

	"github.com/mechta-market/gotemplate/internal/domain/node/model"
)

type Service struct {
	repoDb RepoDbI
}

func New(repoDb RepoDbI) *Service {
	return &Service{repoDb: repoDb}
}

func (s *Service) Create(ctx context.Context, node *model.Node) error {
	if err := s.repoDb.Create(ctx, node); err != nil {
		return fmt.Errorf("repoDb.Create: %w", err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id string) (*model.Node, error) {
	res, err := s.repoDb.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repoDb.Get: %w", err)
	}
	return res, nil
}

func (s *Service) List(ctx context.Context, params model.ListParams) ([]*model.Node, error) {
	res, err := s.repoDb.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("repoDb.List: %w", err)
	}
	return res, nil
}

func (s *Service) Update(ctx context.Context, id string, params model.UpdateParams) error {
	if err := s.repoDb.Update(ctx, id, params); err != nil {
		return fmt.Errorf("repoDb.Update: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repoDb.Delete(ctx, id); err != nil {
		return fmt.Errorf("repoDb.Delete: %w", err)
	}
	return nil
}
