package handlers

import (
	"net/http"
)

func New(handler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", handler.GetAll)
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("GET /users/{id}", handler.GetUserByID)
	mux.HandleFunc("PUT /users/{id}", handler.UpdateUser)
	mux.HandleFunc("/users/{id}", handler.DeleteUser)

	return mux

}
