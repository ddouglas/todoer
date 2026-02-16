package server

import (
	"net/http"

	"github.com/ddouglas/todoer/pkg/types"
)

func (s *Service) handleHome(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	
	// Get filter parameters
	filter := r.URL.Query().Get("filter")
	if filter == "" {
		filter = "all"
	}
	categoryID := r.URL.Query().Get("category")
	
	// Get categories for sidebar
	categories, err := s.categories.Categories(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list categories")
		s.internalServerError(w)
		return
	}
	
	// Get todos based on filter
	var todos []*types.Todo
	if filter == "all" {
		todos, err = s.todos.Todos(ctx)
	} else {
		var catID *string
		if categoryID != "" {
			catID = &categoryID
		}
		todos, err = s.todos.TodosByFilter(ctx, filter, catID)
	}
	
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list todos")
		s.internalServerError(w)
		return
	}
	
	// Find active category if filtering by category
	var activeCategory *types.Category
	if filter == "category" && categoryID != "" {
		for _, cat := range categories {
			if cat.ID == categoryID {
				activeCategory = cat
				break
			}
		}
	}
	
	// Determine selected filter for sidebar highlighting
	selectedFilter := filter
	if filter == "category" && categoryID != "" {
		selectedFilter = categoryID
	}

	data := types.HomePageData{
		Todos:          todos,
		Categories:     categories,
		SelectedFilter: selectedFilter,
		ActiveCategory: activeCategory,
		User:           nil,
	}

	w.Header().Set("Content-Type", "text/html")
	err = s.templates.ExecuteTemplate(w, "page.home", data)
	if err != nil {
		s.logger.WithError(err).Error("failed to render template")
		s.internalServerError(w)
		return
	}

}

func (s *Service) handleNewTodoPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get categories for form
	categories, err := s.categories.Categories(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list categories")
		s.internalServerError(w)
		return
	}

	data := types.TodoDetailPageData{
		Todo:       nil,
		Categories: categories,
		IsNew:      true,
		User:       nil,
	}

	w.Header().Set("Content-Type", "text/html")
	err = s.templates.ExecuteTemplate(w, "page.todo-detail", data)
	if err != nil {
		s.logger.WithError(err).Error("failed to render template")
		s.internalServerError(w)
		return
	}
}

func (s *Service) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	
	// Get the todo
	todo, err := s.todos.Todo(ctx, id)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to get todo")
		http.NotFound(w, r)
		return
	}
	
	// Get categories for form
	categories, err := s.categories.Categories(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list categories")
		s.internalServerError(w)
		return
	}

	data := types.TodoDetailPageData{
		Todo:       todo,
		Categories: categories,
		IsNew:      false,
		User:       nil,
	}

	w.Header().Set("Content-Type", "text/html")
	err = s.templates.ExecuteTemplate(w, "page.todo-detail", data)
	if err != nil {
		s.logger.WithError(err).Error("failed to render template")
		s.internalServerError(w)
		return
	}
}

// Modal handlers - return just the form content for HTMX
func (s *Service) handleNewTodoModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Get categories for form
	categories, err := s.categories.Categories(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list categories")
		s.internalServerError(w)
		return
	}

	data := types.TodoDetailPageData{
		Todo:       nil,
		Categories: categories,
		IsNew:      true,
		User:       nil,
	}

	w.Header().Set("Content-Type", "text/html")
	err = s.templates.ExecuteTemplate(w, "component.modal.todo-detail", data)
	if err != nil {
		s.logger.WithError(err).Error("failed to render template")
		s.internalServerError(w)
		return
	}
}

func (s *Service) handleEditTodoModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	
	// Get the todo
	todo, err := s.todos.Todo(ctx, id)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to get todo")
		http.NotFound(w, r)
		return
	}
	
	// Get categories for form
	categories, err := s.categories.Categories(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list categories")
		s.internalServerError(w)
		return
	}

	data := types.TodoDetailPageData{
		Todo:       todo,
		Categories: categories,
		IsNew:      false,
		User:       nil,
	}

	w.Header().Set("Content-Type", "text/html")
	err = s.templates.ExecuteTemplate(w, "component.modal.todo-detail", data)
	if err != nil {
		s.logger.WithError(err).Error("failed to render template")
		s.internalServerError(w)
		return
	}
}

func (s *Service) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	categories, err := s.categories.Categories(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list categories")
		s.internalServerError(w)
		return
	}

	// Return JSON or template based on Accept header
	// For now, just return as HTMX template
	w.Header().Set("Content-Type", "text/html")
	for _, cat := range categories {
		err = s.templates.ExecuteTemplate(w, "component.category.item", cat)
		if err != nil {
			s.logger.WithError(err).Error("failed to render category")
		}
	}
}

func (s *Service) handleNewCategoryModal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	data := struct {
		IsNew bool
	}{
		IsNew: true,
	}

	w.Header().Set("Content-Type", "text/html")
	err := s.templates.ExecuteTemplate(w, "component.modal.category-form", data)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to render template")
		s.internalServerError(w)
		return
	}
}
