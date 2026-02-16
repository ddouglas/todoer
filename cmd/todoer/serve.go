package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddouglas/todoer/internal/server"
	"github.com/ddouglas/todoer/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var serveCommand = &cli.Command{
	Name:   "serve",
	Usage:  "Start the HTTP Server",
	Action: serve,
}

func serve(cCtx *cli.Context) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{})

	err := LoadConfig(cCtx.String("config"))
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	pool, err := pgxpool.New(ctx, c.Database.URL)
	if err != nil {
		logger.WithError(err).Fatal("failed to open db connection")
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to ping db")
	}

	todoStore := store.NewTodoRepository(pool)
	categoryStore := store.NewCategoryRepository(pool)
	nagStore := store.NewNagRepository(pool)

	srv := server.New(8080, logger, todoStore, categoryStore, nagStore)

	go func() {
		logger.WithField("port", 8080).Info("listening")
		err := srv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("listen err")
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = srv.Stop(shutdownCtx)
	if err != nil {
		logger.WithError(err).Error("server shutdown failed")
		os.Exit(1)
	}

	logger.Info("gracefully shutdown....bye!")
	return nil
}
