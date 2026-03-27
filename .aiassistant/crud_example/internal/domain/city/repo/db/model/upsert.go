package model

import (
	"time"

	"github.com/mechta-market/gotemplate/internal/domain/city/model"
)

type Upsert struct {
	Id string

	ModifiedAt *time.Time
	Active     *bool
	Name       *string
	Dsc        *string
}

func (m *Upsert) CreateColumnMap() map[string]any {
	result := map[string]any{}

	if m.ModifiedAt != nil {
		result["modified_at"] = *m.ModifiedAt
	}
	if m.Active != nil {
		result["active"] = *m.Active
	}
	if m.Name != nil {
		result["name"] = *m.Name
	}
	if m.Dsc != nil {
		result["dsc"] = *m.Dsc
	}

	return result
}

func (m *Upsert) ReturningColumnMap() map[string]any {
	return map[string]any{}
}

func (m *Upsert) UpdateColumnMap() map[string]any {
	return m.CreateColumnMap()
}

func (m *Upsert) PKColumnMap() map[string]any {
	return map[string]any{
		"id": m.Id,
	}
}

// DTO

func DecodeUpsert(v *model.Edit) *Upsert {
	return &Upsert{
		ModifiedAt: v.ModifiedAt,
		Active:     v.Active,
		Name:       v.Name,
		Dsc:        v.Dsc,
	}
}
