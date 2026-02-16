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

const nagTableName = "todo_nag_settings"

var nagTableColumns = utils.StructTagValues(&types.NagSettings{})

type NagRepository struct {
	pool *pgxpool.Pool
}

func NewNagRepository(pool *pgxpool.Pool) *NagRepository {
	return &NagRepository{pool: pool}
}

func (r *NagRepository) NagSettingsByTodoID(ctx context.Context, todoID string) (*types.NagSettings, error) {
	query, args, err := psql().
		Select(nagTableColumns...).
		From(nagTableName).
		Where(sq.Eq{"todo_id": todoID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query: %w", err)
	}

	var nag = new(types.NagSettings)

	err = pgxscan.Get(ctx, r.pool, nag, query, args...)
	if err != nil {
		return nil, err
	}

	return nag, nil
}

func (r *NagRepository) CreateNag(ctx context.Context, nag *types.NagSettings) error {
	now := time.Now()
	nag.ID = utils.NanoID()
	nag.CreatedAt = now
	nag.UpdatedAt = now

	query, args, err := psql().Insert(nagTableName).SetMap(utils.StructToMap(nag)).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate create nag query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec create nag query: %w", err)
	}

	return nil
}

func (r *NagRepository) UpdateNag(ctx context.Context, id string, nag *types.NagSettings) error {
	nag.ID = id
	nag.UpdatedAt = time.Now()

	query, args, err := psql().Update(nagTableName).SetMap(utils.StructToMap(nag)).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate update nag query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec update nag query: %w", err)
	}

	return nil
}

func (r *NagRepository) DeleteNag(ctx context.Context, todoID string) error {
	query, args, err := psql().Delete(nagTableName).Where(sq.Eq{"todo_id": todoID}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to generate delete nag query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec delete nag query: %w", err)
	}

	return nil
}

func (r *NagRepository) ActiveNags(ctx context.Context) ([]*types.NagSettings, error) {
	// Query to find all nag settings that should send a notification now
	// Start from nag_settings (smaller table), join to todos for status/due_date checks
	query := `
		SELECT n.*
		FROM todo_nag_settings n
		INNER JOIN todos t ON n.todo_id = t.id
		WHERE 
			n.enabled = true
			AND t.status = $1
			AND t.due_date IS NOT NULL
			AND t.due_date <= NOW()
			AND t.due_date + (n.duration_minutes || ' minutes')::interval >= NOW()
			AND (
				n.last_nagged_at IS NULL 
				OR n.last_nagged_at + (n.interval_minutes || ' minutes')::interval <= NOW()
			)
		ORDER BY t.due_date ASC
	`

	var nags []*types.NagSettings
	err := pgxscan.Select(ctx, r.pool, &nags, query, types.StatusNotStarted)
	if err != nil {
		return nil, fmt.Errorf("failed to query active nags: %w", err)
	}

	return nags, nil
}
