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

	mux := http.NewServeMux()

	handler := middleware.Logging(
		middleware.Auth(mux))

	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("/users/{id}", userHandler.DeleteUser)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
