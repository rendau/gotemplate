package model

import domainModel "github.com/mechta-market/gotemplate/internal/domain/node/model"

type Upsert struct {
	ID   string
	Name *string
	Ord  *int
}

func (m *Upsert) CreateColumnMap() map[string]any {
	result := map[string]any{
		"id": m.ID,
	}
	if m.Name != nil {
		result["name"] = *m.Name
	}
	if m.Ord != nil {
		result["ord"] = *m.Ord
	}
	return result
}

func (m *Upsert) UpdateColumnMap() map[string]any {
	result := map[string]any{}
	if m.Name != nil {
		result["name"] = *m.Name
	}
	if m.Ord != nil {
		result["ord"] = *m.Ord
	}
	return result
}

// DTO

func DecodeCreate(v *domainModel.Node) *Upsert {
	return &Upsert{
		ID:   v.ID,
		Name: &v.Name,
		Ord:  &v.Ord,
	}
}

func DecodeUpdate(id string, v domainModel.UpdateParams) *Upsert {
	return &Upsert{
		ID:   id,
		Name: v.Name,
		Ord:  v.Ord,
	}
}
