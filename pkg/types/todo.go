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
	Priority    string     `db:"priority"`
	CategoryID  *string    `db:"category_id"`
	DueDate     *time.Time `db:"due_date"`
	ReminderAt  *time.Time `db:"reminder_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	
	// Joined fields (not in DB)
	Category *Category `db:"-"`
}

const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)
