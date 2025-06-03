package example

import (
	"context"
	"github.com/mechta-market/gotemplate/internal/domain/example/model"
)

type iService interface {
	Get(ctx context.Context, id string) (*model.Main, error)
}
