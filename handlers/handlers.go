package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fnuritdinov/firstService/models"
	"github.com/fnuritdinov/firstService/pkg/logger"
	"github.com/fnuritdinov/firstService/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
	logger  logger.Logger
}

func NewUserHandler(service *service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := h.service.GetUsers()
	if err != nil {
		h.logger.Error("error from service", zap.String("method", "GetUsers"))
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
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.Create(user)
	if err != nil {
		h.logger.Error("error from service",
			zap.String("method", "Create"))
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Update(id, user.Name)
	if err != nil {
		h.logger.Error("error from service",
			zap.String("method", "Update"))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id := models.StrToInt(idStr)

	err := h.service.Delete(id)
	h.logger.Error("error from service",
		zap.String("method", "Delete"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
