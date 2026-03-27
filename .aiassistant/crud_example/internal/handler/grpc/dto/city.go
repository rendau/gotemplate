package dto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mechta-market/gotemplate/internal/domain/city/model"
	"github.com/mechta-market/gotemplate/pkg/proto/example_v1"
)

func EncodeCityMain(v *model.Main, _ int) *example_v1.CityMain {
	if v == nil {
		return nil
	}

	return &example_v1.CityMain{
		Id:         v.Id,
		CreatedAt:  timestamppb.New(v.CreatedAt),
		ModifiedAt: timestamppb.New(v.ModifiedAt),
		Active:     v.Active,
		Name:       v.Name,
		Dsc:        v.Dsc,
	}
}

func DecodeCityListReq(v *example_v1.CityListReq) *model.ListReq {
	if v == nil {
		return nil
	}

	return &model.ListReq{
		ListParams: DecodeListParams(v.ListParams),
		Active:     v.Active,
	}
}

func DecodeCityCreateReq(v *example_v1.CityCreateReq) *model.Edit {
	if v == nil {
		return nil
	}

	return &model.Edit{
		Active: v.Active,
		Name:   v.Name,
		Dsc:    v.Dsc,
	}
}

func DecodeCityUpdateReq(v *example_v1.CityUpdateReq) *model.Edit {
	if v == nil {
		return nil
	}

	return &model.Edit{
		Active: v.Active,
		Name:   v.Name,
		Dsc:    v.Dsc,
	}
}
