package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/ShadyZiedan/gophermart/internal/config"
	"github.com/ShadyZiedan/gophermart/internal/handlers"
	"github.com/ShadyZiedan/gophermart/internal/infrastructure"
	"github.com/ShadyZiedan/gophermart/internal/infrastructure/repositories"
	"github.com/ShadyZiedan/gophermart/internal/integration"
	"github.com/ShadyZiedan/gophermart/internal/logger"
	"github.com/ShadyZiedan/gophermart/internal/security"
	"github.com/ShadyZiedan/gophermart/internal/server"
	"github.com/ShadyZiedan/gophermart/internal/services"
)

func main() {
	err := logger.Init("info")
	if err != nil {
		panic(err)
	}
	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Log.Error("Error parsing config", zap.Error(err))
		os.Exit(1)
	}

	if err := infrastructure.Migrate(cfg.DatabaseURI); err != nil {
		logger.Log.Error("Error migrating database", zap.Error(err))
		os.Exit(1)
	}
	pgxPool, err := pgxpool.New(context.Background(), cfg.DatabaseURI)
	if err != nil {
		logger.Log.Error("Error initializing database connection", zap.Error(err))
		os.Exit(1)
	}
	secretKey, err := security.GenerateRandomKey(32)
	if err != nil {
		logger.Log.Error("Error generating random key", zap.Error(err))
		os.Exit(1)
	}

	authService := services.NewAuthService(
		secretKey,
		repositories.NewUsersRepository(pgxPool),
	)
	orderService := services.NewOrderService(repositories.NewOrdersRepository(pgxPool))
	balanceService := services.NewBalanceService(repositories.NewBalanceRepository(pgxPool))
	accrualService := integration.NewAccrualService(cfg.AccrualSystemAddress, orderService)
	serverHandler := handlers.NewRouter(authService, orderService, balanceService, accrualService)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Info("starting HTTP server", zap.String("address", cfg.RunAddress))
		err := server.NewWebServer(cfg.RunAddress, serverHandler).Run(ctx)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Log.Info("gracefully shutting down server")
			} else {
				logger.Log.Error("error starting web server", zap.Error(err))
			}
		}
	}()
	wg.Wait()
}
