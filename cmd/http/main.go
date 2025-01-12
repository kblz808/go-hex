package main

import (
	"context"
	"go-hex/internal/adapter/config"
	"go-hex/internal/adapter/handler/http"
	"go-hex/internal/adapter/storage/postgres"
	"go-hex/internal/adapter/storage/postgres/repository"
	"go-hex/internal/core/service"
	"log/slog"
	"os"

	"github.com/docker/docker/volume/service"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// logger.Set(config.App)
	slog.Info("Starting app", "app", config.App.Name, "env", config.App.Env)

	ctx := context.Background()
	db, err := postgres.New(ctx, *config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("successfully connected to the database", "db", config.DB.Connection)

	err = db.Migrate()
	if err != nil {
		slog.Error("error migrating database", "error", err)
		os.Exit(1)
	}

	slog.Info("Successfully migrated the database")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	// router, err := http.NewRoute
}
