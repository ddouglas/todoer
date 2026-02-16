package types

import "time"

type Category struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Color     string    `db:"color"`
	Icon      *string   `db:"icon"`
	SortOrder int       `db:"sort_order"`
	CreatedAt time.Time `db:"created_at"`
}

type Todo struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	Completed   bool       `db:"completed"`
	Status      string     `db:"status"`
	Priority    string     `db:"priority"`
	CategoryID  *string    `db:"category_id"`
	DueDate     *time.Time `db:"due_date"`
	ReminderAt  *time.Time `db:"reminder_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	
	// Joined fields (not in DB)
	Category    *Category    `db:"-"`
	NagSettings *NagSettings `db:"-"`
}

type NagSettings struct {
	ID              string     `db:"id"`
	TodoID          string     `db:"todo_id"`
	Enabled         bool       `db:"enabled"`
	IntervalMinutes int        `db:"interval_minutes"`
	DurationMinutes int        `db:"duration_minutes"`
	LastNaggedAt    *time.Time `db:"last_nagged_at"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"

	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)
