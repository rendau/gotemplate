package model

import (
	"time"

	"github.com/mechta-market/gotemplate/internal/domain/city/model"
)

type Select struct {
	Id         string
	CreatedAt  time.Time
	ModifiedAt time.Time
	Active     bool
	Name       string
	Dsc        string
}

func (m *Select) ListColumnMap() map[string]any {
	return map[string]any{
		"id":          &m.Id,
		"created_at":  &m.CreatedAt,
		"modified_at": &m.ModifiedAt,
		"active":      &m.Active,
		"name":        &m.Name,
		"dsc":         &m.Dsc,
	}
}

func (m *Select) PKColumnMap() map[string]any {
	return map[string]any{
		"id": m.Id,
	}
}

func (m *Select) DefaultSortColumns() []string {
	return []string{
		"active desc",
		"name",
	}
}

// DTO

func EncodeSelect(v *Select, _ int) *model.Main {
	return &model.Main{
		Id:         v.Id,
		CreatedAt:  v.CreatedAt,
		ModifiedAt: v.ModifiedAt,
		Active:     v.Active,
		Name:       v.Name,
		Dsc:        v.Dsc,
	}
}
