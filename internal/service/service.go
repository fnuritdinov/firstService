package service

import (
	"context"
	"fmt"

	"github.com/fnuritdinov/firstService/internal/models"
	"github.com/fnuritdinov/firstService/internal/service/eventLogs"
	"github.com/fnuritdinov/firstService/internal/storage"
	"github.com/fnuritdinov/firstService/pkg/errors"
	"github.com/fnuritdinov/firstService/pkg/utils"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(storage *storage.UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

var Users = map[string]string{
	"user1": "Ali",
	"user2": "Vali",
	"user3": "Anton",
}

const admin = "admin"

func (u UserService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := u.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from u.storage.GetAll: %w", err)
	}

	err = eventLogs.Audit(admin, "Получил список пользователей", eventLogs.GetAll)
	if err != nil {
		return nil, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	return users, nil
}

func (u UserService) GetUserByID(id int) (models.User, error) {
	err := utils.ValidateID(id)
	if err != nil {
		return models.User{}, errors.ErrorFromValidateID
	}

	user, err := u.storage.GetByID(id)
	if err != nil {
		return models.User{}, err
	}

	err = eventLogs.Audit(admin, "Получил пользователя по ID", eventLogs.GetUserByID)
	if err != nil {
		return models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}
	return user, nil
}

func (u UserService) Create(ctx context.Context, user models.User) error {
	err := user.ValidateID()
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = user.ValidateStrEmpty()
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = u.storage.Create(ctx, user)
	if err != nil {
		return errors.ErrorNotFound
	}

	err = eventLogs.Audit(admin, "Создал пользователя", eventLogs.Create)
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	return nil

}

func (u UserService) Update(id int, updatedName string) error {
	err := utils.ValidateID(id)
	if err != nil {
		return errors.ErrorFromValidateID
	}

	err = utils.ValidateStrEmpty(updatedName)
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = u.storage.Update(id, updatedName)
	if err != nil {
		return errors.ErrorNotFound
	}

	err = eventLogs.Audit(admin, "Изменил пользователя", eventLogs.Update)
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}
	return nil
}

func (u UserService) Delete(id int) error {
	err := utils.ValidateID(id)
	if err != nil {
		return errors.ErrorFromValidateID
	}

	err = u.storage.Delete(id)
	if err != nil {
		return errors.ErrorNotFound
	}

	err = eventLogs.Audit(admin, "Удалил пользователя", eventLogs.Delete)
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}
	return nil
}
