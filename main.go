package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fnuritdinov/firstService/handlers"
	"github.com/fnuritdinov/firstService/internal/config"
	"github.com/fnuritdinov/firstService/internal/consumer"
	"github.com/fnuritdinov/firstService/internal/rate_limiter"
	"github.com/fnuritdinov/firstService/internal/service"
	"github.com/fnuritdinov/firstService/internal/service/eventBus"
	"github.com/fnuritdinov/firstService/internal/storage"
	"github.com/fnuritdinov/firstService/middleware"
	"github.com/fnuritdinov/firstService/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	logger, err := logger.New(true)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}
	bus := eventBus.NewBus(10)

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	consumer.StartAuditConsumer(ctx, &wg, bus, *logger)

	rate := rate_limiter.New()
	go rate.WorkerClear(ctx)

	cfg, err := config.New("./config/config.env")
	if err != nil {
		log.Fatal("config.New", err)
	}

	userStorage := storage.NewUserStorage(cfg.Storage)
	userService := service.NewUserService(userStorage, bus)
	userHandler := handlers.NewUserHandler(userService, *logger)

	handler := handlers.New(userHandler)

	handler2 := middleware.RateLimit(rate, middleware.Logging(
		middleware.Auth(handler)))

	server := &http.Server{
		Addr:    cfg.HttpPort,
		Handler: handler2,
	}

	go func() {
		logger.Info("server started")

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutdown started")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("shutdown error", zap.Error(err))
	}

	wg.Wait()
	logger.Info("shutdown completed")

}
