package handlers

import (
	"net/http"
)

func New(handler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("PUT /users/{id}", handler.UpdateUser)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)
	mux.HandleFunc("GET /users/all", handler.GetAll)
	mux.HandleFunc("GET /users", handler.Get)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("PATCH /users/{id}/password", handler.UpdatePassword)
	mux.HandleFunc("PATCH /users/{id}/age", handler.UpdateAge)
	mux.HandleFunc("GET /users/stats", handler.GetUsersStats)

	return mux
}
