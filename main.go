package main

import (
	http "net/http"

	"github.com/fnuritdinov/firstService/handlers"
	"github.com/fnuritdinov/firstService/middleware"
	"github.com/fnuritdinov/firstService/service"
	"github.com/fnuritdinov/firstService/storage"
)

func main() {

	userStorage := storage.NewUserStorage("data/user.json")
	userService := service.NewUserService(userStorage)
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()

	handler := middleware.Logging(
		middleware.Auth(mux))

	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
	mux.HandleFunc("/users/{id}", userHandler.DeleteUser)

	http.ListenAndServe(":8080", handler)
}
