package db

import (
	"strings"

	"github.com/mechta-market/gotemplate/internal/domain/node/model"
)

var (
	allowedSortFields = map[string]string{
		"id":   "id",
		"name": "name",
		"ord":  "ord",
	}
)

func (r *Repo) getConditions(_ model.ListParams) (map[string]any, map[string][]any) {
	conditions := make(map[string]any)
	conditionExps := make(map[string][]any)

	return conditions, conditionExps
}

func constructSort(sort []string) []string {
	if len(sort) == 0 {
		return []string{"ord asc", "id asc"}
	}

	res := make([]string, 0, len(sort))
	for _, s := range sort {
		direction := "asc"
		field := s
		if strings.HasPrefix(s, "-") {
			direction = "desc"
			field = strings.TrimPrefix(s, "-")
		}
		col, ok := allowedSortFields[field]
		if !ok {
			continue
		}
		res = append(res, col+" "+direction)
	}
	if len(res) == 0 {
		return []string{"ord asc", "id asc"}
	}
	return res
}
