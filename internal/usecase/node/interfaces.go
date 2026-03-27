package node

import (
	"context"

	"github.com/mechta-market/gotemplate/internal/domain/node/model"
)

type ServiceI interface {
	Create(ctx context.Context, node *model.Node) error
	Get(ctx context.Context, id string) (*model.Node, error)
	List(ctx context.Context, params model.ListParams) ([]*model.Node, error)
	Update(ctx context.Context, id string, params model.UpdateParams) error
	Delete(ctx context.Context, id string) error
}
