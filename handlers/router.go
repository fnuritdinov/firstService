package handlers

import (
	"net/http"
)

func New(handler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users/all", handler.GetAll)
	mux.HandleFunc("GET /users", handler.Get)
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)
	mux.HandleFunc("PUT /users/{id}", handler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", handler.DeleteUser)
	mux.HandleFunc("POST /login", handler.Login)

	return mux
}
