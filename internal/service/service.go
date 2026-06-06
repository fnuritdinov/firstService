package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/fnuritdinov/firstService/internal/models"
	"github.com/fnuritdinov/firstService/internal/service/eventBus"
	"github.com/fnuritdinov/firstService/internal/service/eventLogs"
	"github.com/fnuritdinov/firstService/internal/storage"
	"github.com/fnuritdinov/firstService/pkg/errors"
	"github.com/fnuritdinov/firstService/pkg/utils"
)

type UserService struct {
	storage   *storage.UserStorage
	eventsBus *eventBus.Bus
}

func NewUserService(storage *storage.UserStorage, eventsBus *eventBus.Bus) *UserService {
	return &UserService{
		storage:   storage,
		eventsBus: eventsBus,
	}
}

var Users = map[string]string{
	"user1": "Ali",
	"user2": "Vali",
	"user3": "Anton",
}

const admin = "admin"

func (u UserService) Login(ctx context.Context, user models.User) error {
	err := utils.ValidateStrEmpty(user.Name)
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = utils.ValidateStrEmpty(user.Password)
	if err != nil {
		return errors.ErrorFromValidateStrEmpty
	}

	err = u.storage.Login(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u UserService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := u.storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from u.storage.GetAll: %w", err)
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.GetAll, "Получил всех пользователей")
	if err != nil {
		return []models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	u.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.GetAll,
		UserID: rand.Int(),
	})

	return users, nil
}

func (u UserService) Get(ctx context.Context) ([]models.User, error) {
	users, err := u.storage.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from u.storage.GetAll: %w", err)
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.Get, "Получил активных пользователей")
	if err != nil {
		return []models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	u.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Get,
		UserID: rand.Int(),
	})

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

	err = eventLogs.Audit(rand.Int(), eventLogs.GetUserByID, "Получил пользователя по ID")
	if err != nil {
		return models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	u.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.GetUserByID,
		UserID: rand.Int(),
	})

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

	err = eventLogs.Audit(rand.Int(), eventLogs.Create, "Создал пользователя")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	u.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Create,
		UserID: rand.Int(),
	})

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

	err = eventLogs.Audit(rand.Int(), eventLogs.Update, "Изменил пользователя")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	u.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Update,
		UserID: rand.Int(),
	})

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

	err = eventLogs.Audit(rand.Int(), eventLogs.Delete, "Удалил пользователя")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	u.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Delete,
		UserID: rand.Int(),
	})

	return nil
}
