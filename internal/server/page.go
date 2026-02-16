package server

import (
	"net/http"

	"github.com/ddouglas/todoer/pkg/types"
)

type homePageData struct {
	Todos []*types.Todo
	User  *struct{}
}

func (s *Service) handleHome(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	todos, err := s.todos.Todos(ctx)
	if err != nil {
		s.logger.WithError(err).WithContext(ctx).Error("failed to list todos")
		s.internalServerError(w)
		return
	}

	data := homePageData{
		Todos: todos,
		User:  nil,
	}

	w.Header().Set("Content-Type", "text/html")
	err = s.templates.ExecuteTemplate(w, "page.home", data)
	if err != nil {
		s.logger.WithError(err).Error("failed to render template")
		s.internalServerError(w)
		return
	}

}
