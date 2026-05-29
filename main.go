package main

import (
	http "net/http"

	"firtService/handlers"
	"firtService/middleware"
	"firtService/storage"
)

func main() {

	st := &storage.UserStorage{
		FileName: "data/users.json",
	}

	h := &handlers.UserHandler{
		Storage: st,
	}

	mux := http.NewServeMux()

	handler := middleware.Logging(
		middleware.Auth(mux))

	mux.HandleFunc("GET /users", h.GetUsers)
	mux.HandleFunc("POST /users", h.CreateUser)
	mux.HandleFunc("GET /users/{ID}", h.GetUserByID)
	mux.HandleFunc("PUT /users/{ID}", h.UpdateUser)
	mux.HandleFunc("/users/{id}", h.DeleteUser)

	http.ListenAndServe(":8080", handler)
}
