package city

import (
	"context"
	"fmt"

	"github.com/mechta-market/gotemplate/internal/domain/city/model"

	"github.com/mechta-market/gotemplate/internal/domain/common/util"
	"github.com/mechta-market/gotemplate/internal/errs"
)

type Usecase struct {
	svc ServiceI
}

func New(svc ServiceI) *Usecase {
	return &Usecase{
		svc: svc,
	}
}

func (u *Usecase) List(ctx context.Context, pars *model.ListReq) ([]*model.Main, int64, error) {
	if err := util.RequirePageSize(pars.ListParams, 0); err != nil {
		return nil, 0, err
	}

	items, tCount, err := u.svc.List(ctx, pars)
	if err != nil {
		return nil, 0, fmt.Errorf("svc.List: %w", err)
	}

	return items, tCount, err
}

func (u *Usecase) Create(ctx context.Context, obj *model.Edit) error {
	err := u.validateEdit(obj, true)
	if err != nil {
		return err
	}

	err = u.svc.Create(ctx, obj)
	if err != nil {
		return fmt.Errorf("svc.Create: %w", err)
	}

	return nil
}

func (u *Usecase) Get(ctx context.Context, id string) (*model.Main, error) {
	result, _, err := u.svc.Get(ctx, id, true)
	if err != nil {
		return nil, fmt.Errorf("svc.Get: %w", err)
	}

	return result, nil
}

func (u *Usecase) Update(ctx context.Context, id string, obj *model.Edit) error {
	if id == "" {
		return errs.IdRequired
	}
	err := u.validateEdit(obj, false)
	if err != nil {
		return err
	}

	err = u.svc.Update(ctx, id, obj)
	if err != nil {
		return fmt.Errorf("svc.Update: %w", err)
	}

	return nil
}

func (u *Usecase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errs.IdRequired
	}
	err := u.svc.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("svc.Delete: %w", err)
	}

	return nil
}

func (u *Usecase) validateEdit(obj *model.Edit, forCreate bool) error {
	return nil
}
