package db

import "github.com/mechta-market/gotemplate/internal/domain/city/model"

var (
	allowedSortFields = map[string]string{
		"id":          "id",
		"created_at":  "created_at",
		"modified_at": "modified_at",
		"name":        "name",
	}
)

func (r *Repo) getConditions(pars *model.ListReq) (map[string]any, map[string][]any) {
	conditions := make(map[string]any)
	conditionExps := make(map[string][]any)

	if pars == nil {
		return conditions, conditionExps
	}

	return conditions, conditionExps
}
