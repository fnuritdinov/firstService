package service

import (
	"fmt"

	"github.com/fnuritdinov/firstService/models"
	"github.com/fnuritdinov/firstService/storage"
)

type UserService struct {
	storage *storage.UserStorage
}

func NewUserService(storage *storage.UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (u UserService) GetUsers() ([]models.User, error) {
	users, err := u.storage.GetAll()
	if err != nil {
		fmt.Println("storage error ", err)
		return nil, err
	}
	return users, nil
}

func (u UserService) GetUserByID(id int) (models.User, error) {
	err := models.ValidateID(id)
	if err != nil {
		return models.User{}, models.ErrorFromValidateID
	}

	user, err := u.storage.GetByID(id)
	if err != nil {
		fmt.Println("Ошибка в бд GetByID")
		return models.User{}, models.ErrorNotFound
	}
	return user, nil
}

func (u UserService) Create(user models.User) error {
	err := models.ValidateID(user.ID)
	if err != nil {
		return models.ErrorFromValidateID
	}

	err = models.ValidateStrEmpty(user.Name)
	if err != nil {
		return models.ErrorFromValidateStrEmpty
	}

	err = u.storage.Create(user)
	if err != nil {
		fmt.Println("Ошибка в бд Create")
		return models.ErrorNotFound
	}
	return nil

}

func (u UserService) Update(id int, name string) error {
	err := models.ValidateID(id)
	if err != nil {
		return models.ErrorFromValidateID
	}

	err = models.ValidateStrEmpty(name)
	if err != nil {
		return models.ErrorFromValidateStrEmpty
	}

	err = u.storage.Update(id, name)
	if err != nil {
		fmt.Println("Ошибка в бд Update")
		return models.ErrorNotFound
	}
	return nil
}

func (u UserService) Delete(id int) error {
	err := models.ValidateID(id)
	if err != nil {
		return models.ErrorFromValidateID
	}

	err = u.storage.Delete(id)
	if err != nil {
		fmt.Println("Ошибка в бд Delete")
		return models.ErrorNotFound
	}
	return nil
}
