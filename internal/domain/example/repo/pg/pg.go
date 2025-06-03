package pg

import (
	"context"
	"github.com/mechta-market/gotemplate/internal/domain/example/model"
)

type pg struct{}

func (r *pg) Get(ctx context.Context, id string) (*model.Main, error) {
	return &model.Main{}, nil
}
