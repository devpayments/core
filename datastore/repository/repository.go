package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Model[E any] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}

type Repository[M Model[E], E any] interface {
	FindByID(ctx context.Context, id string) (E, error)
}

type BaseRepository[M Model[E], E any] struct {
	db        *sqlx.DB
	tableName string
}

func NewBaseRepository[M Model[E], E any](db *sqlx.DB, tableName string) *BaseRepository[M, E] {
	return &BaseRepository[M, E]{
		db:        db,
		tableName: tableName,
	}
}

func (r *BaseRepository[M, E]) Create(ctx context.Context, entity *E) (int64, error) {
	var m M
	model := m.FromEntity(*entity).(M)
	insertSQL, _, err := goqu.Insert(r.tableName).Rows(
		model,
	).ToSQL()

	if err != nil {
		panic(err)
	}

	res, err := r.db.ExecContext(ctx, insertSQL)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *BaseRepository[M, E]) FindByID(ctx context.Context, id uuid.UUID) (*E, error) {
	var model M

	selectSql, _, err := goqu.Select("*").From(r.tableName).Where(goqu.C("id").Eq(id.String())).ToSQL()
	if err != nil {
		panic(err)
	}
	row := r.db.QueryRowxContext(ctx, selectSql)
	err = row.StructScan(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	e := model.ToEntity()
	return &e, nil
}

func (r *BaseRepository[M, E]) FindOne(ctx context.Context, whereMap map[string]any) (*E, error) {
	where := goqu.Ex{}
	for index, element := range whereMap {
		where[index] = element
	}
	selectSql, _, err := goqu.
		Select("*").
		From(r.tableName).
		Where(
			where,
		).
		ToSQL()
	if err != nil {
		panic(err)
	}

	row := r.db.QueryRowxContext(ctx, selectSql)
	var model M
	err = row.StructScan(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	e := model.ToEntity()
	return &e, nil
}

func (r *BaseRepository[M, E]) FindAll(ctx context.Context, whereMap map[string]any) (*E, error) {
	where := goqu.Ex{}
	for index, element := range whereMap {
		where[index] = element
	}
	selectSql, _, err := goqu.
		Select("*").
		From(r.tableName).
		Where(
			where,
		).
		ToSQL()
	if err != nil {
		panic(err)
	}

	row, err := r.db.QueryxContext(ctx, selectSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var model M
	err = row.StructScan(&model)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	e := model.ToEntity()
	return &e, nil
}
