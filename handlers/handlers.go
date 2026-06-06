package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fnuritdinov/firstService/internal/models"
	"github.com/fnuritdinov/firstService/internal/service"
	errs "github.com/fnuritdinov/firstService/pkg/errors"
	"github.com/fnuritdinov/firstService/pkg/logger"
	"github.com/fnuritdinov/firstService/pkg/utils"
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

type userRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type loginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("error from json.NewDecoder", zap.Error(err))
		http.Error(w, "error from json.NewDecocder", http.StatusBadRequest)
		return
	}

	err = h.service.Login(r.Context(), models.User{
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, errs.ErrorFromValidateStrEmpty) {
			h.logger.Error("error from validateStrEmpty", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	users, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("error from h.service.GetAll", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.Get(r.Context())
	if err != nil {
		h.logger.Error("error from h.service.GetAll", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id := utils.StrToInt(idStr)

	user, err := h.service.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrorFromValidateID) {
			h.logger.Error("error from validate",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user userRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Create(r.Context(), models.User{
		ID:       user.ID,
		Name:     user.Name,
		IsActive: true,
	})
	if err != nil {
		if errors.Is(err, errs.ErrorFromValidateStrEmpty) {
			h.logger.Error("error from errs.ErrorFromValidateStrEmpty",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, errs.ErrorFromValidateID) {
			h.logger.Error("error from ValidateID",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusCreated)

}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id := utils.StrToInt(idStr)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.Update(id, user.Name)
	if err != nil {
		if errors.Is(err, errs.ErrorFromValidateID) {
			h.logger.Error("error from ValidateID", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, errs.ErrorFromValidateStrEmpty) {
			h.logger.Error("error from ValidateStrEmpty", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.logger.Error("error from service",
			zap.Error(err))
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id := utils.StrToInt(idStr)

	err := h.service.Delete(id)
	h.logger.Error("error from service",
		zap.Error(err))
	if err != nil {
		if errors.Is(err, errs.ErrorFromValidateID) {
			h.logger.Error("error from ValidateID", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
