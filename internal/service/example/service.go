package example

import (
	"context"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Routine(ctx context.Context) error {
	return nil
}
