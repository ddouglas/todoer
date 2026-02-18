package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ddouglas/todoer/internal/notify"
	"github.com/ddouglas/todoer/internal/store"
	"github.com/ddouglas/todoer/pkg/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var workerCommand = &cli.Command{
	Name:  "worker",
	Usage: "Run the notification worker to send nag reminders",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to config file",
			EnvVars: []string{"TODOER_CONFIG"},
			Value:   "config.yml",
		},
		&cli.DurationFlag{
			Name:    "interval",
			Aliases: []string{"i"},
			Usage:   "How often to check for nags",
			Value:   time.Minute,
		},
	},
	Action: runWorker,
}

func runWorker(cCtx *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{})

	// Load config
	err := LoadConfig(cCtx.String("config"))
	if err != nil {
		logger.WithError(err).Fatal("failed to load config")
	}

	// Validate NTFY config
	if c.NTFY.URL == "" || c.NTFY.Topic == "" {
		logger.Fatal("NTFY configuration required (url and topic)")
	}

	// Connect to database
	pool, err := pgxpool.New(ctx, c.Database.URL)
	if err != nil {
		logger.WithError(err).Fatal("failed to open db connection")
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		logger.WithError(err).Fatal("failed to ping db")
	}

	// Initialize store and notify client
	nagStore := store.NewNagRepository(pool)
	todoStore := store.NewTodoRepository(pool)
	ntfyClient := notify.NewNTFYClient(c.NTFY.URL, c.NTFY.Topic)

	logger.WithFields(logrus.Fields{
		"interval":   cCtx.Duration("interval"),
		"ntfy_url":   c.NTFY.URL,
		"ntfy_topic": c.NTFY.Topic,
	}).Info("starting notification worker")

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create ticker
	ticker := time.NewTicker(cCtx.Duration("interval"))
	defer ticker.Stop()

	// Run immediately on start
	processNags(ctx, logger, nagStore, todoStore, ntfyClient)

	// Main loop
	for {
		select {
		case <-ticker.C:
			processNags(ctx, logger, nagStore, todoStore, ntfyClient)
		case sig := <-sigChan:
			logger.WithField("signal", sig).Info("received shutdown signal")
			return nil
		case <-ctx.Done():
			logger.Info("context cancelled, shutting down")
			return nil
		}
	}
}

func processNags(ctx context.Context, logger *logrus.Logger, nagStore *store.NagRepository, todoStore *store.TodoRepository, ntfyClient *notify.NTFYClient) {
	nags, err := nagStore.ActiveNags(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to fetch active nags")
		return
	}

	if len(nags) == 0 {
		logger.Debug("no active nags to process")
		return
	}

	logger.WithField("count", len(nags)).Info("processing nags")

	// Collect all todo IDs
	todoIDs := make([]string, len(nags))
	for i, nag := range nags {
		todoIDs[i] = nag.TodoID
	}

	// Fetch all todos at once
	todos, err := todoStore.TodosByIDs(ctx, todoIDs)
	if err != nil {
		logger.WithError(err).Error("failed to fetch todos")
		return
	}

	// Map todos by ID for quick lookup
	todoMap := make(map[string]*types.Todo)
	for _, todo := range todos {
		todoMap[todo.ID] = todo
	}

	// Process each nag
	for _, nag := range nags {
		todo, ok := todoMap[nag.TodoID]
		if !ok {
			logger.WithField("todo_id", nag.TodoID).Warn("todo not found for nag")
			continue
		}
		nag.Todo = todo

		err = sendNagNotification(ctx, logger, ntfyClient, nag)
		if err != nil {
			logger.WithError(err).WithField("todo_id", nag.TodoID).Error("failed to send nag notification")
			continue
		}

		// Update last_nagged_at timestamp
		now := time.Now()
		nag.LastNaggedAt = &now
		err = nagStore.UpdateNag(ctx, nag.ID, nag)
		if err != nil {
			logger.WithError(err).WithField("nag_id", nag.ID).Error("failed to update last_nagged_at")
		}

		logger.WithFields(logrus.Fields{
			"todo_id":    nag.TodoID,
			"todo_title": nag.Todo.Title,
		}).Info("sent nag notification")
	}
}

func sendNagNotification(ctx context.Context, logger *logrus.Logger, ntfyClient *notify.NTFYClient, nag *types.NagSettings) error {
	if nag.Todo == nil {
		return fmt.Errorf("nag has no associated todo")
	}

	title := "⏰ Todo Reminder"
	message := fmt.Sprintf("Don't forget: %s", nag.Todo.Title)

	// Add urgency based on priority
	priority := 3 // default
	tags := []string{"alarm_clock"}

	if nag.Todo.Priority == types.PriorityHigh {
		priority = 5
		tags = append(tags, "fire")
		title = "🔥 Urgent Todo Reminder"
	} else if nag.Todo.Priority == types.PriorityMedium {
		priority = 4
	}

	// Add description if available
	if nag.Todo.Description != nil && *nag.Todo.Description != "" {
		message = fmt.Sprintf("%s\n\n%s", message, *nag.Todo.Description)
	}

	return ntfyClient.Send(ctx, &notify.Message{
		Title:    title,
		Message:  message,
		Priority: priority,
		Tags:     tags,
	})
}
