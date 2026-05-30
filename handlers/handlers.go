package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fnuritdinov/firstService/models"
	"github.com/fnuritdinov/firstService/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.service.GetUsers()
	if err != nil {
		fmt.Println("Ошибка в сервисе handler/GetUsers")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id := models.StrToInt(idStr)

	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		fmt.Println("Ошибка в сервисе handler/GetUserByID")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("encode error:", err)
		return
	}

}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Ошибка при парсинге hanlder/CreateUser")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.Create(user)
	if err != nil {
		fmt.Println("Ошибка в сервисе handler/CreateUser")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user created",
	})

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id := models.StrToInt(idStr)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Ошибка при парсинге handler/UpdateUser")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Update(id, user.Name)
	if err != nil {
		fmt.Println("Ошибка в сервисе handler/UpdateUser")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id := models.StrToInt(idStr)

	err := h.service.Delete(id)
	fmt.Println("Ошибка в сервисе handler/DeleteUser")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
