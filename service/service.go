package service

import (
	"github.com/fnuritdinov/firstService/models"
	"github.com/fnuritdinov/firstService/pkg/errors"
	"github.com/fnuritdinov/firstService/pkg/utils"
	"github.com/fnuritdinov/firstService/storage"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(storage *storage.UserStorage) *UserService {
	return &UserService{
		storage: storage,
	}
}

func (u UserService) GetUsers() ([]models.User, error) {
	users, err := u.storage.GetAll()
	if err != nil {
		return nil, err
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
	return user, nil
}

func (u UserService) Create(user models.User) error {
	err := user.ValidateID()
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = user.ValidateStrEmpty()
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = u.storage.Create(user)
	if err != nil {
		return errors.ErrorNotFound
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
	return nil
}
