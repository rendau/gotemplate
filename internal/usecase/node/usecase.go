package node

import (
	"context"

	domainmodel "github.com/mechta-market/gotemplate/internal/domain/node/model"
)

type Usecase struct {
	service ServiceI
}

func New(service ServiceI) *Usecase {
	return &Usecase{service: service}
}

func (u *Usecase) Create(ctx context.Context, in CreateIn) error {
	return u.service.Create(ctx, &domainmodel.Node{
		ID:   in.ID,
		Name: in.Name,
		Ord:  in.Ord,
	})
}

func (u *Usecase) Get(ctx context.Context, id string) (*Node, error) {
	item, err := u.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Node{Node: *item}, nil
}

func (u *Usecase) List(ctx context.Context, in ListIn) ([]*Node, error) {
	list, err := u.service.List(ctx, domainmodel.ListParams{
		ListParams: in.ListParams,
	})
	if err != nil {
		return nil, err
	}

	res := make([]*Node, 0, len(list))
	for _, item := range list {
		res = append(res, &Node{Node: *item})
	}

	return res, nil
}

func (u *Usecase) Update(ctx context.Context, id string, in UpdateIn) error {
	return u.service.Update(ctx, id, domainmodel.UpdateParams{
		Name: in.Name,
		Ord:  in.Ord,
	})
}

func (u *Usecase) Delete(ctx context.Context, id string) error {
	return u.service.Delete(ctx, id)
}
