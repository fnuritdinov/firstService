package service

import (
	"context"
	"fmt"

	"github.com/fnuritdinov/firstService/internal/models"
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

func (u UserService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := u.storage.GetAll(ctx)
	if err != nil {
		fmt.Errorf("error from u.storage.GetAll %w", err)
		return nil, err
	}
	return users, nil
}

func (u UserService) GetUserByID(id int) (models.User, error) {
	err := utils.ValidateID(id)
	if err != nil {
		fmt.Errorf("error from utils.ValidateID %w", err)
		return models.User{}, errors.ErrorFromValidateID
	}

	user, err := u.storage.GetByID(id)
	if err != nil {
		fmt.Errorf("error from u.storage.GetByID %w", err)
		return models.User{}, err
	}
	return user, nil
}

func (u UserService) Create(ctx context.Context, user models.User) error {
	err := user.ValidateID()
	if err != nil {
		fmt.Errorf("error from user.ValidateID %w", err)
		return errors.ErrorFromValidateStrEmpty
	}

	err = user.ValidateStrEmpty()
	if err != nil {
		fmt.Errorf("error from user.ValidateStrEmpty %w", err)
		return errors.ErrorFromValidateStrEmpty
	}

	err = u.storage.Create(ctx, user)
	if err != nil {
		fmt.Errorf("error from u.storage.Create %w", err)
		return errors.ErrorNotFound
	}
	return nil

}

func (u UserService) Update(id int, updatedName string) error {
	err := utils.ValidateID(id)
	if err != nil {
		fmt.Errorf("error from utils.ValidateID %w", err)
		return errors.ErrorFromValidateID
	}

	err = utils.ValidateStrEmpty(updatedName)
	if err != nil {
		fmt.Errorf("error from utils.ValidateStrEmpty %w", err)
		return errors.ErrorFromValidateStrEmpty
	}

	err = u.storage.Update(id, updatedName)
	if err != nil {
		fmt.Errorf("error from u.storage.Update %w", err)
		return errors.ErrorNotFound
	}
	return nil
}

func (u UserService) Delete(id int) error {
	err := utils.ValidateID(id)
	if err != nil {
		fmt.Errorf("error from utils.ValidateID %w", err)
		return errors.ErrorFromValidateID
	}

	err = u.storage.Delete(id)
	if err != nil {
		fmt.Errorf("error from u.storage.Delete %w", err)
		return errors.ErrorNotFound
	}
	return nil
}
