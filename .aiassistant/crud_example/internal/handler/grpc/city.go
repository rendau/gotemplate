package grpc

import (
	"context"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/mechta-market/gotemplate/internal/handler/grpc/dto"
	usecase "github.com/mechta-market/gotemplate/internal/usecase/city"
	"github.com/mechta-market/gotemplate/pkg/proto/example_v1"
	"github.com/mechta-market/gotemplate/pkg/proto/common"
)

type City struct {
	example_v1.UnsafeCityServer
	usecase *usecase.Usecase
}

func NewCity(usecase *usecase.Usecase) *City {
	return &City{
		usecase: usecase,
	}
}

func (h *City) List(ctx context.Context, req *example_v1.CityListReq) (*example_v1.CityListRep, error) {
	if req.ListParams == nil {
		req.ListParams = &common.ListParamsSt{}
	}

	items, tCount, err := h.usecase.List(ctx, dto.DecodeCityListReq(req))
	if err != nil {
		return nil, err
	}

	return &example_v1.CityListRep{
		PaginationInfo: &common.PaginationInfoSt{
			Page:       req.ListParams.Page,
			PageSize:   req.ListParams.PageSize,
			TotalCount: tCount,
		},
		Results: lo.Map(items, dto.EncodeCityMain),
	}, nil
}

func (h *City) Get(ctx context.Context, req *example_v1.CityGetReq) (*example_v1.CityMain, error) {
	item, err := h.usecase.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return dto.EncodeCityMain(item, 0), nil
}

func (h *City) Create(ctx context.Context, req *example_v1.CityCreateReq) (*emptypb.Empty, error) {
	err := h.usecase.Create(ctx, dto.DecodeCityCreateReq(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *City) Update(ctx context.Context, req *example_v1.CityUpdateReq) (*emptypb.Empty, error) {
	err := h.usecase.Update(ctx, req.PkId, dto.DecodeCityUpdateReq(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *City) Delete(ctx context.Context, req *example_v1.CityGetReq) (*emptypb.Empty, error) {
	err := h.usecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
