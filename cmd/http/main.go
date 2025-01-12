package main

import (
	"context"
	"fmt"
	"go-hex/internal/adapter/auth/paseto"
	"go-hex/internal/adapter/config"
	"go-hex/internal/adapter/handler/http"
	"go-hex/internal/adapter/storage/postgres"
	"go-hex/internal/adapter/storage/postgres/repository"
	"go-hex/internal/core/service"
	"log/slog"
	"os"
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

	token, err := paseto.New(config.Token)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, token)
	authHandler := http.NewAuthHandler(authService)

	router, err := http.NewRouter(config.HTTP, token, *userHandler, *authHandler)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the http server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting http server", "error", err)
		os.Exit(1)
	}
}
