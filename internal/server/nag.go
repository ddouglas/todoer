package server

import (
	"net/http"
	"strconv"

	"github.com/ddouglas/todoer/pkg/types"
)

func (s *Service) handleGetNagSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todoID := r.PathValue("id")

	nag, err := s.nags.NagSettingsByTodoID(ctx, todoID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch nag settings")
		s.internalServerError(w)
		return
	}

	// For now, just return success - template rendering will be added later
	err = s.templates.ExecuteTemplate(w, "component.nag.settings", nag)
	if err != nil {
		s.logger.WithError(err).WithField("template_name", "component.nag.settings").Error("failed to render template")
		s.internalServerError(w)
		return
	}
}

func (s *Service) handlePostNagSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todoID := r.PathValue("id")

	err := r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	enabled := r.FormValue("enabled") == "true"
	intervalStr := r.FormValue("interval_minutes")
	durationStr := r.FormValue("duration_minutes")

	nag := &types.NagSettings{
		TodoID:  todoID,
		Enabled: enabled,
	}

	// Parse interval_minutes (default to 10)
	if intervalStr != "" {
		interval, err := strconv.Atoi(intervalStr)
		if err == nil {
			nag.IntervalMinutes = interval
		} else {
			nag.IntervalMinutes = 10
		}
	} else {
		nag.IntervalMinutes = 10
	}

	// Parse duration_minutes (default to 60)
	if durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err == nil {
			nag.DurationMinutes = duration
		} else {
			nag.DurationMinutes = 60
		}
	} else {
		nag.DurationMinutes = 60
	}

	err = s.nags.CreateNag(ctx, nag)
	if err != nil {
		s.logger.WithError(err).Error("failed to create nag settings")
		s.internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) handlePutNagSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todoID := r.PathValue("id")

	// Fetch existing nag settings
	existing, err := s.nags.NagSettingsByTodoID(ctx, todoID)
	if err != nil {
		s.logger.WithError(err).Error("failed to fetch existing nag settings")
		s.internalServerError(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		s.logger.WithError(err).Error("failed to parse http request form")
		s.internalServerError(w)
		return
	}

	enabled := r.FormValue("enabled") == "true"
	intervalStr := r.FormValue("interval_minutes")
	durationStr := r.FormValue("duration_minutes")

	nag := &types.NagSettings{
		ID:      existing.ID,
		TodoID:  todoID,
		Enabled: enabled,
	}

	// Parse interval_minutes
	if intervalStr != "" {
		interval, err := strconv.Atoi(intervalStr)
		if err == nil {
			nag.IntervalMinutes = interval
		} else {
			nag.IntervalMinutes = existing.IntervalMinutes
		}
	} else {
		nag.IntervalMinutes = existing.IntervalMinutes
	}

	// Parse duration_minutes
	if durationStr != "" {
		duration, err := strconv.Atoi(durationStr)
		if err == nil {
			nag.DurationMinutes = duration
		} else {
			nag.DurationMinutes = existing.DurationMinutes
		}
	} else {
		nag.DurationMinutes = existing.DurationMinutes
	}

	err = s.nags.UpdateNag(ctx, existing.ID, nag)
	if err != nil {
		s.logger.WithError(err).Error("failed to update nag settings")
		s.internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) handleDeleteNagSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	todoID := r.PathValue("id")

	err := s.nags.DeleteNag(ctx, todoID)
	if err != nil {
		s.logger.WithError(err).Error("failed to delete nag settings")
		s.internalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
