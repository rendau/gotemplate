package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"

	commonRepoPg "github.com/mechta-market/gotemplate/internal/domain/common/repo/pg"
	"github.com/mechta-market/gotemplate/internal/domain/node/model"
	repoModel "github.com/mechta-market/gotemplate/internal/domain/node/repo/db/model"
)

type Repo struct {
	*commonRepoPg.Base
}

func New(con *pgxpool.Pool) *Repo {
	return &Repo{Base: commonRepoPg.NewBase(con)}
}

func (r *Repo) Create(ctx context.Context, node *model.Node) error {
	m := repoModel.DecodeCreate(node)
	columns := m.CreateColumnMap()
	qb := r.QB.Insert("node").SetMap(columns)

	q, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build create node query: %w", err)
	}

	if _, err = r.Con.Exec(ctx, q, args...); err != nil {
		return fmt.Errorf("exec create node query: %w", err)
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, id string) (*model.Node, error) {
	q, args, err := r.QB.Select("id", "name", "ord").
		From("node").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get node query: %w", err)
	}

	res := &repoModel.Select{}
	if err = r.Con.QueryRow(ctx, q, args...).Scan(&res.ID, &res.Name, &res.Ord); err != nil {
		return nil, fmt.Errorf("exec get node query: %w", err)
	}

	return repoModel.EncodeSelect(res, 0), nil
}

func (r *Repo) List(ctx context.Context, params model.ListParams) ([]*model.Node, error) {
	conditions, _ := r.getConditions(params)
	sort := constructSort(params.Sort)

	q := r.QB.Select("id", "name", "ord").From("node")
	for k, v := range conditions {
		q = q.Where(k, v)
	}
	if len(sort) > 0 {
		q = q.OrderBy(sort...)
	}
	if params.PageSize > 0 {
		q = q.Limit(uint64(params.PageSize))
		if params.Page > 1 {
			q = q.Offset(uint64((params.Page - 1) * params.PageSize))
		}
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list node query: %w", err)
	}

	rows, err := r.Con.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec list node query: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*repoModel.Select, error) {
		m := &repoModel.Select{}
		if scanErr := row.Scan(&m.ID, &m.Name, &m.Ord); scanErr != nil {
			return nil, fmt.Errorf("scan node row: %w", scanErr)
		}
		return m, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect list node rows: %w", err)
	}

	return lo.Map(items, repoModel.EncodeSelect), nil
}

func (r *Repo) Update(ctx context.Context, id string, params model.UpdateParams) error {
	m := repoModel.DecodeUpdate(id, params)
	columns := m.UpdateColumnMap()
	if len(columns) == 0 {
		return nil
	}

	qb := r.QB.Update("node").Where("id = ?", m.ID)
	for col, value := range columns {
		qb = qb.Set(col, value)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build update node query: %w", err)
	}

	if _, err = r.Con.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("exec update node query: %w", err)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	sql, args, err := r.QB.Delete("node").Where("id = ?", strings.TrimSpace(id)).ToSql()
	if err != nil {
		return fmt.Errorf("build delete node query: %w", err)
	}

	if _, err = r.Con.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("exec delete node query: %w", err)
	}

	return nil
}
