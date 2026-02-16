package types

// Page data structs for templates
type HomePageData struct {
	Todos          []*Todo
	Categories     []*Category
	SelectedFilter string // "all", "today", "week", or category ID
	ActiveCategory *Category
	User           *struct{} // Placeholder for future auth
}

type TodoDetailPageData struct {
	Todo       *Todo
	Categories []*Category
	IsNew      bool
	User       *struct{} // Placeholder for future auth
}
