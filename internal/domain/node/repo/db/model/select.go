package model

import domainModel "github.com/mechta-market/gotemplate/internal/domain/node/model"

type Select struct {
	ID   string
	Name string
	Ord  int
}

// DTO

func EncodeSelect(v *Select, _ int) *domainModel.Node {
	return &domainModel.Node{
		ID:   v.ID,
		Name: v.Name,
		Ord:  v.Ord,
	}
}
