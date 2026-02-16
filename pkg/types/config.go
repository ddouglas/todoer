package types

import (
	_ "github.com/kelseyhightower/envconfig"
)

type Config struct {
	Environment string `envconfig:"ENVIRONMENT" required:"true"`
	Port        uint   `envconfig:"PORT" required:"true"`

	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

// Page data structs for templates
type HomePageData struct {
	Todos           []*Todo
	Categories      []*Category
	SelectedFilter  string // "all", "today", "week", or category ID
	ActiveCategory  *Category
}

type TodoDetailPageData struct {
	Todo       *Todo
	Categories []*Category
	IsNew      bool
}
