package model

import commonModel "github.com/mechta-market/gotemplate/internal/domain/common/model"

type Node struct {
	ID   string
	Name string
	Ord  int
}

type ListParams struct {
	commonModel.ListParams
}

type UpdateParams struct {
	Name *string
	Ord  *int
}
