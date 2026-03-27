package model

import (
	"time"

	commonModel "github.com/mechta-market/gotemplate/internal/domain/common/model"
)

type Main struct {
	Id         string
	CreatedAt  time.Time
	ModifiedAt time.Time
	Active     bool
	Name       string
	Dsc        string
}

type ListReq struct {
	commonModel.ListParams

	Active *bool
}

type Edit struct {
	Id         *string
	ModifiedAt *time.Time
	Active     *bool
	Name       *string
	Dsc        *string
}
