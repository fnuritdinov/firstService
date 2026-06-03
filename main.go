package main

import (
	"log"
	http "net/http"

	"github.com/fnuritdinov/firstService/handlers"
	"github.com/fnuritdinov/firstService/internal/service"
	"github.com/fnuritdinov/firstService/internal/storage"
	"github.com/fnuritdinov/firstService/middleware"
	"github.com/fnuritdinov/firstService/pkg/logger"
)

func main() {
	logger, err := logger.New(true)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}
	userStorage := storage.NewUserStorage("data/user.json")
	userService := service.NewUserService(userStorage)
	userHandler := handlers.NewUserHandler(userService, *logger)

	handler := handlers.New(userHandler)

	handler2 := middleware.Logging(
		middleware.Auth(handler))

	err = http.ListenAndServe(":8080", handler2)
	if err != nil {
		log.Fatal(err)
	}
}
