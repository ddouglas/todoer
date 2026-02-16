package store

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/ddouglas/todoer/internal/utils"
	"github.com/ddouglas/todoer/pkg/types"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

const categoryTableName = "categories"

var categoryTableColumns = utils.StructTagValues(&types.Category{})

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) Category(ctx context.Context, id string) (*types.Category, error) {

	query, args, err := psql().
		Select(categoryTableColumns...).
		From(categoryTableName).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var category = new(types.Category)

	err = pgxscan.Get(ctx, r.pool, category, query, args...)

	return category, err

}

func (r *CategoryRepository) Categories(ctx context.Context) ([]*types.Category, error) {

	query, args, err := psql().
		Select(categoryTableColumns...).
		From(categoryTableName).
		OrderBy("sort_order ASC", "name ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var categories = make([]*types.Category, 0)

	err = pgxscan.Select(ctx, r.pool, &categories, query, args...)

	return categories, err

}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *types.Category) error {

	category.ID = utils.NanoID()
	category.CreatedAt = time.Now()

	query, args, err := psql().Insert(categoryTableName).SetMap(utils.StructToMap(category)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate create category query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec create category query: %w", err)
	}

	return nil

}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, id string, category *types.Category) error {

	category.ID = id

	query, args, err := psql().Update(categoryTableName).SetMap(utils.StructToMap(category)).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate update category query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec update category query: %w", err)
	}

	return nil

}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {

	query, args, err := psql().Delete(categoryTableName).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate delete category query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec delete category query: %w", err)
	}

	return nil

}
