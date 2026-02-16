package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/flow"
	"github.com/ddouglas/todoer/internal/store"
	"github.com/sirupsen/logrus"
)

//go:embed templates
var templateFS embed.FS

type Service struct {
	port   uint
	logger *logrus.Logger

	todos *store.TodoRepository

	templates *template.Template
	server    *http.Server
}

const defaultTimeout = time.Second * 5

func New(
	port uint,
	logger *logrus.Logger,

	todos *store.TodoRepository,
) *Service {

	mux := flow.New()

	s := &Service{
		port:   port,
		logger: logger,

		todos: todos,

		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", port),
			Handler:           mux,
			ReadTimeout:       defaultTimeout,
			ReadHeaderTimeout: defaultTimeout,
			WriteTimeout:      defaultTimeout,
			MaxHeaderBytes:    512,
		},
	}

	s.buildRouter(mux)
	s.loadTemplates()

	return s

}

func (s *Service) Start() error {
	return s.server.ListenAndServe()
}

func (s *Service) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Service) buildRouter(r *flow.Mux) {

	r.Use(s.LoggingMiddleware)

	r.HandleFunc("/", s.handleHome, http.MethodGet)
	r.HandleFunc("/todos", s.handlePostTodo, http.MethodPost)
	r.HandleFunc("/todos/:id", s.handlePatchTodo, http.MethodPatch)
	r.HandleFunc("/todos/:id", s.handleDeleteTodo, http.MethodDelete)

}

func (s *Service) loadTemplates() {

	tmpl := template.New("")
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		data, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return fmt.Errorf("failed to read template at %s: %w", path, err)
		}

		_, err = tmpl.Parse(string(data))
		if err != nil {
			return fmt.Errorf("failed to parse template at %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		s.logger.WithError(err).Fatal("failed to load templates")
	}

	s.templates = tmpl

}

func (s *Service) internalServerError(w http.ResponseWriter) {
	http.Error(w, "interal server error", http.StatusInternalServerError)
}
