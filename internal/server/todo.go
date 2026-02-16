package server

import (
	"net/http"
	"time"

	"github.com/ddouglas/todoer/pkg/types"
)

func (s *Service) handlePostTodo(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")
	priority := r.FormValue("priority")
	categoryID := r.FormValue("category_id")
	dueDateStr := r.FormValue("due_date")
	reminderStr := r.FormValue("reminder_at")

	todo := &types.Todo{
		Title:    title,
		Priority: priority,
	}

	// Optional fields
	if description != "" {
		todo.Description = &description
	}
	if categoryID != "" {
		todo.CategoryID = &categoryID
	}
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			todo.DueDate = &dueDate
		}
	}
	if reminderStr != "" {
		reminder, err := time.Parse("2006-01-02T15:04", reminderStr)
		if err == nil {
			todo.ReminderAt = &reminder
		}
	}

	err = s.todos.CreateTodo(ctx, todo)
	if err != nil {
		s.logger.WithError(err).Error("failed to create todo")
		s.internalServerError(w)
		return
	}

	// Reload todo with category
	todo, err = s.todos.Todo(ctx, todo.ID)
	if err != nil {
		s.logger.WithError(err).Error("failed to reload todo")
	}

	err = s.templates.ExecuteTemplate(w, "component.todo.item", todo)
	if err != nil {
		s.logger.WithError(err).WithField("template_name", "component.todo.item").Error("failed to render template")
		s.internalServerError(w)
		return
	}

	count, err := s.todos.TodoCount(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch updated list of todos after a create")
		return
	}

	if count <= 1 {
		err = s.templates.ExecuteTemplate(w, "component.todo.empty-delete", nil)
		if err != nil {
			s.logger.WithError(err).Error("failed to render empty-state delete")
		}
	}

}

func (s *Service) handlePutTodo(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	id := r.PathValue("id")

	title := r.FormValue("title")
	description := r.FormValue("description")
	priority := r.FormValue("priority")
	categoryID := r.FormValue("category_id")
	dueDateStr := r.FormValue("due_date")
	reminderStr := r.FormValue("reminder_at")
	completed := r.FormValue("completed") == "on"

	todo := &types.Todo{
		Title:     title,
		Priority:  priority,
		Completed: completed,
	}

	// Optional fields
	if description != "" {
		todo.Description = &description
	}
	if categoryID != "" {
		todo.CategoryID = &categoryID
	}
	if dueDateStr != "" {
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err == nil {
			todo.DueDate = &dueDate
		}
	}
	if reminderStr != "" {
		reminder, err := time.Parse("2006-01-02T15:04", reminderStr)
		if err == nil {
			todo.ReminderAt = &reminder
		}
	}

	err = s.todos.UpdateTodo(ctx, id, todo)
	if err != nil {
		s.logger.WithError(err).Error("failed to update todo")
		s.internalServerError(w)
		return
	}

	// Reload todo with category
	todo, err = s.todos.Todo(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("failed to reload todo")
	}

	err = s.templates.ExecuteTemplate(w, "component.todo.item", todo)
	if err != nil {
		s.logger.WithError(err).WithField("template_name", "component.todo.item").Error("failed to render template")
		s.internalServerError(w)
		return
	}

}

func (s *Service) handlePatchTodo(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	id := r.PathValue("id")

	todo, err := s.todos.Todo(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch todo or todo is unknown")
		s.internalServerError(w)
		return
	}

	todo.Completed = r.FormValue("completed") == "on"

	err = s.todos.UpdateTodo(ctx, id, todo)
	if err != nil {
		s.logger.WithError(err).Error("failed to update todo")
		s.internalServerError(w)
		return
	}

	err = s.templates.ExecuteTemplate(w, "component.todo.item", todo)
	if err != nil {
		s.logger.WithError(err).WithField("template_name", "component.todo.item").Error("failed to render template")
		s.internalServerError(w)
		return
	}

}

func (s *Service) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := r.PathValue("id")

	err := s.todos.DeleteTodo(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete todo")
		s.internalServerError(w)
		return
	}

	count, err := s.todos.TodoCount(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch updated list of todos after a delete")
		return
	}

	if count == 0 {
		w.Header().Set("Content-Type", "text/html")
		err = s.templates.ExecuteTemplate(w, "component.todo.empty-oob", nil)
		if err != nil {
			s.logger.WithError(err).Error("failed to write oob swap for empty-state id")
		}
	}

}
