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
	service service.IUserService
	logger  logger.Logger
}

func NewUserHandler(service service.IUserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

type userRequest struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	IsActive bool   `json:"isActive"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   int    `json:"userID"`
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		u.logger.Error("error from json.NewDecoder", zap.Error(err))
		http.Error(w, "error from json.NewDecocder", http.StatusBadRequest)
		return
	}

	err = u.service.Login(r.Context(), models.Auth{
		Login:    req.Login,
		Password: req.Password,
		UserID:   req.UserID,
	})
	if err != nil {
		if errors.Is(err, errs.ErrFromValidate) {
			u.logger.Error("error from validateStrEmpty", zap.Error(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user userRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.service.Create(r.Context(), models.User{
		Name:     user.Name,
		Age:      user.Age,
		IsActive: true,
	})
	if err != nil {
		if errors.Is(err, errs.ErrFromValidate) {
			u.logger.Error("error from errs.ErrorFromValidateStrEmpty",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, errs.ErrFromValidateID) {
			u.logger.Error("error from ValidateID",
				zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	id, err := utils.StrToInt(r.PathValue("id"))
	if err != nil {
		u.logger.Error("error from utils.StrToInt")
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.service.Update(r.Context(), id, user.Name)
	if err != nil {
		if errors.Is(err, errs.ErrFromValidateID) {
			u.logger.Error("error from ValidateID", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, errs.ErrFromValidate) {
			u.logger.Error("error from ValidateStrEmpty", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.logger.Error("error from service",
			zap.Error(err))
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	users, err := u.service.GetAll(r.Context())
	if err != nil {
		u.logger.Error("error from h.service.GetAll", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	users, err := u.service.Get(r.Context())
	if err != nil {
		u.logger.Error("error from h.service.GetAll", zap.Error(err))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := utils.StrToInt(idStr)
	if err != nil {
		u.logger.Error("error from utils.StrToInt")
		return
	}

	user, err := u.service.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrFromValidateID) {
			u.logger.Error("error from validate",
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

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := utils.StrToInt(r.PathValue("id"))
	if err != nil {
		u.logger.Error("error from utils.StrToInt")
		return
	}

	err = u.service.Delete(r.Context(), id)
	u.logger.Error("error from service",
		zap.Error(err))
	if err != nil {
		if errors.Is(err, errs.ErrFromValidateID) {
			u.logger.Error("error from ValidateID", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type updatePasswordReq struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (u *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

	id, err := utils.StrToInt(r.PathValue("id"))
	if err != nil {
		u.logger.Error("error from utils.StrToInt")
		return
	}

	var pass updatePasswordReq
	err = json.NewDecoder(r.Body).Decode(&pass)
	if err != nil {
		u.logger.Error("error from json.NewDecoder")
		return
	}

	u.logger.Info("handler UpdatePassword",
		zap.Int("user_id", id),
		zap.Any("request", pass),
	)

	err = u.service.UpdatePassword(r.Context(), id, models.Auth{
		OldPassword: pass.OldPassword,
		NewPassword: pass.NewPassword,
	})
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, errs.ErrBadRequest) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, errs.ErrFromValidate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		if errors.Is(err, errs.ErrWrongPassword) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		u.logger.Error("error from u.service.UpdatePassword",
			zap.Error(err))
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "password updated successfully",
	})
}

func (u *UserHandler) UpdateAge(w http.ResponseWriter, r *http.Request) {

	id, err := utils.StrToInt(r.PathValue("id"))
	if err != nil {
		u.logger.Error("error from utils.StrToInt")
		return
	}

	var req userRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = u.service.UpdateAge(r.Context(), id, models.User{
		Age: req.Age,
	})
	if err != nil {
		if errors.Is(err, errs.ErrFromValidate) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, errs.ErrIsActiveFalse) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		u.logger.Error("error from u.service.UpdateAge",
			zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "age updated successfully",
	})
}

func (u *UserHandler) GetUsersStats(w http.ResponseWriter, r *http.Request) {

	users, err := u.service.GetUsersStats(r.Context())
	if err != nil {
		u.logger.Error("error from u.service.GetUsersStats",
			zap.Error(err))
		http.Error(w, "invalid error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		u.logger.Error("error from json.Encoder", zap.Error(err))
	}
}
