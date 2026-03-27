package node

import (
	commonModel "github.com/mechta-market/gotemplate/internal/domain/common/model"
	domainmodel "github.com/mechta-market/gotemplate/internal/domain/node/model"
)

type Node struct {
	domainmodel.Node
}

type CreateIn struct {
	ID   string
	Name string
	Ord  int
}

type ListIn struct {
	commonModel.ListParams
}

type UpdateIn struct {
	Name *string
	Ord  *int
}
