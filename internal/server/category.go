package server

import (
	"net/http"

	"github.com/ddouglas/todoer/pkg/types"
)

func (s *Service) handlePostCategory(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	name := r.FormValue("name")
	color := r.FormValue("color")
	icon := r.FormValue("icon")

	category := &types.Category{
		Name:  name,
		Color: color,
	}

	if icon != "" {
		category.Icon = &icon
	}

	err = s.categories.CreateCategory(ctx, category)
	if err != nil {
		s.logger.WithError(err).Error("failed to create category")
		s.internalServerError(w)
		return
	}

	// Redirect to home on success
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) handlePutCategory(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	id := r.PathValue("id")

	name := r.FormValue("name")
	color := r.FormValue("color")
	icon := r.FormValue("icon")

	category := &types.Category{
		Name:  name,
		Color: color,
	}

	if icon != "" {
		category.Icon = &icon
	}

	err = s.categories.UpdateCategory(ctx, id, category)
	if err != nil {
		s.logger.WithError(err).Error("failed to update category")
		s.internalServerError(w)
		return
	}

	// Redirect to home on success
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Service) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	id := r.PathValue("id")

	err := s.categories.DeleteCategory(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete category")
		s.internalServerError(w)
		return
	}

	// Redirect to home on success
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
