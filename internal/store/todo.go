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

const todoTableName = "todos"

var todoTableColumns = utils.StructTagValues(&types.Todo{})

type TodoRepository struct {
	pool *pgxpool.Pool
}

func NewTodoRepository(pool *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{pool: pool}
}

func (r *TodoRepository) Todo(ctx context.Context, id string) (*types.Todo, error) {

	query, args, err := psql().
		Select(todoTableColumns...).
		From(todoTableName).
		Where(sq.Eq{"todos.id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var todo = new(types.Todo)

	err = pgxscan.Get(ctx, r.pool, todo, query, args...)
	if err != nil {
		return nil, err
	}

	// Load category if present
	if todo.CategoryID != nil {
		category, err := r.loadCategory(ctx, *todo.CategoryID)
		if err == nil {
			todo.Category = category
		}
	}

	return todo, nil

}

func (r *TodoRepository) Todos(ctx context.Context) ([]*types.Todo, error) {

	query, args, err := psql().
		Select(todoTableColumns...).
		From(todoTableName).
		OrderBy("completed ASC", "created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var todos = make([]*types.Todo, 0)

	err = pgxscan.Select(ctx, r.pool, &todos, query, args...)
	if err != nil {
		return nil, err
	}

	// Load categories for all todos
	for _, todo := range todos {
		if todo.CategoryID != nil {
			category, err := r.loadCategory(ctx, *todo.CategoryID)
			if err == nil {
				todo.Category = category
			}
		}
	}

	return todos, nil

}

func (r *TodoRepository) TodosByFilter(ctx context.Context, filter string, categoryID *string) ([]*types.Todo, error) {
	builder := psql().
		Select(todoTableColumns...).
		From(todoTableName)

	now := time.Now()

	switch filter {
	case "today":
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		builder = builder.Where(sq.And{
			sq.GtOrEq{"due_date": startOfDay},
			sq.Lt{"due_date": endOfDay},
		})
	case "week":
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
		startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
		endOfWeek := startOfWeek.AddDate(0, 0, 7)
		builder = builder.Where(sq.And{
			sq.GtOrEq{"due_date": startOfWeek},
			sq.Lt{"due_date": endOfWeek},
		})
	case "category":
		if categoryID != nil {
			builder = builder.Where(sq.Eq{"category_id": *categoryID})
		}
	}

	builder = builder.OrderBy("completed ASC", "due_date ASC NULLS LAST", "created_at DESC")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var todos = make([]*types.Todo, 0)

	err = pgxscan.Select(ctx, r.pool, &todos, query, args...)
	if err != nil {
		return nil, err
	}

	// Load categories for all todos
	for _, todo := range todos {
		if todo.CategoryID != nil {
			category, err := r.loadCategory(ctx, *todo.CategoryID)
			if err == nil {
				todo.Category = category
			}
		}
	}

	return todos, nil
}

func (r *TodoRepository) TodoCount(ctx context.Context) (int, error) {

	query, args, err := psql().Select("COUNT(*)").From(todoTableName).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to generate query: %w", err)
	}

	var count int

	err = pgxscan.Get(ctx, r.pool, &count, query, args...)

	return count, err

}

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *types.Todo) error {

	now := time.Now()
	todo.ID = utils.NanoID()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	query, args, err := psql().Insert(todoTableName).SetMap(utils.StructToMap(todo)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate create todo query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec create todo query: %w", err)
	}

	return nil

}

func (r *TodoRepository) UpdateTodo(ctx context.Context, id string, todo *types.Todo) error {

	todo.ID = id
	todo.UpdatedAt = time.Now()

	query, args, err := psql().Update(todoTableName).SetMap(utils.StructToMap(todo)).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate update todo query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec update todo query: %w", err)
	}

	return nil

}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id string) error {

	query, args, err := psql().Delete(todoTableName).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate update todo query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec update todo query: %w", err)
	}

	return nil

}

// loadCategory is a helper method to load a category by ID
func (r *TodoRepository) loadCategory(ctx context.Context, id string) (*types.Category, error) {
	categoryColumns := utils.StructTagValues(&types.Category{})
	
	query, args, err := psql().
		Select(categoryColumns...).
		From("categories").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate category query: %w", err)
	}

	var category = new(types.Category)
	err = pgxscan.Get(ctx, r.pool, category, query, args...)
	
	return category, err
}
