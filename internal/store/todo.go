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

	query, args, err := psql().Select(todoTableColumns...).From(todoTableName).Where(sq.Eq{
		"id": id,
	}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var todo = new(types.Todo)

	err = pgxscan.Get(ctx, r.pool, todo, query, args...)

	return todo, err

}

func (r *TodoRepository) Todos(ctx context.Context) ([]*types.Todo, error) {

	query, args, err := psql().Select(todoTableColumns...).From(todoTableName).OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var todos = make([]*types.Todo, 0)

	err = pgxscan.Select(ctx, r.pool, &todos, query, args...)

	return todos, err

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

	todo.ID = utils.NanoID()
	todo.CreatedAt = time.Now()

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
