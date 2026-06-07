package util

import (
	"{{ cookiecutter.__go_module }}/internal/domain/common/model"
	"{{ cookiecutter.__go_module }}/internal/errs"
)

const (
	defaultMaxPageSize int64 = 100
)

func RequirePageSize(pars model.ListParams, maxPageSize int64) error {
	if maxPageSize == 0 {
		maxPageSize = defaultMaxPageSize
	}

	if pars.PageSize == 0 || pars.PageSize > maxPageSize {
		return errs.IncorrectPageSize
	}

	return nil
}
