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

func (u UserService) Login(ctx context.Context, user models.User) error {
	err := user.Validate()
	if err != nil {
		return errors.ErrFromValidate
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

func (u UserService) GetUserByID(ctx context.Context, id int) (models.User, error) {
	err := utils.ValidateID(id)
	if err != nil {
		return models.User{}, errors.ErrFromValidateID
	}

	user, err := u.storage.GetByID(ctx, id)
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
	err := user.Validate()
	if err != nil {
		return errors.ErrFromValidate
	}

	err = u.storage.Create(ctx, user)
	if err != nil {
		return errors.ErrNotFound
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

func (u UserService) Update(ctx context.Context, id int, updatedName string) error {
	err := utils.ValidateID(id)
	if err != nil {
		return errors.ErrFromValidateID
	}

	err = utils.ValidateStrEmpty(updatedName)
	if err != nil {
		return errors.ErrFromValidate
	}

	err = u.storage.Update(ctx, id, updatedName)
	if err != nil {
		return errors.ErrNotFound
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

func (u UserService) Delete(ctx context.Context, id int) error {
	err := utils.ValidateID(id)
	if err != nil {
		return errors.ErrFromValidateID
	}

	err = u.storage.Delete(ctx, id)
	if err != nil {
		return errors.ErrNotFound
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
