package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/fnuritdinov/firstService/internal/models"
	"github.com/fnuritdinov/firstService/internal/repository"
	"github.com/fnuritdinov/firstService/internal/service/eventBus"
	"github.com/fnuritdinov/firstService/internal/service/eventLogs"
	"github.com/fnuritdinov/firstService/pkg/errors"
	"github.com/fnuritdinov/firstService/pkg/utils"
)

type IUserService interface {
	Login(ctx context.Context, auth models.Auth) error
	GetAll(ctx context.Context) ([]models.User, error)
	Get(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	Create(ctx context.Context, user models.User) error
	Update(ctx context.Context, id int, updatedName string) error
	Delete(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, id int, pass models.Auth) error
	UpdateAge(ctx context.Context, id int, user models.User) error
}

type service struct {
	repository repository.IUserRepo
	eventsBus  *eventBus.Bus
}

func NewUserService(repository repository.IUserRepo, eventsBus *eventBus.Bus) IUserService {
	return &service{
		repository: repository,
		eventsBus:  eventsBus,
	}
}

func (s *service) Login(ctx context.Context, auth models.Auth) error {
	err := auth.ValidateAuth()
	if err != nil {
		return err
	}

	err = s.repository.Login(ctx, auth)
	if err != nil {
		return err
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.Login, "пользователь залгогинился")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Login,
		UserID: rand.Int(),
	})
	return nil
}

func (s *service) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from u.storage.GetAll: %w", err)
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.GetAll, "Получил всех пользователей")
	if err != nil {
		return []models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.GetAll,
		UserID: rand.Int(),
	})

	return users, nil
}

func (s *service) Get(ctx context.Context) ([]models.User, error) {
	users, err := s.repository.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error from u.storage.GetAll: %w", err)
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.Get, "Получил активных пользователей")
	if err != nil {
		return []models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Get,
		UserID: rand.Int(),
	})

	return users, nil
}

func (s *service) GetUserByID(ctx context.Context, id int) (models.User, error) {
	err := utils.ValidateID(id)
	if err != nil {
		return models.User{}, errors.ErrFromValidateID
	}

	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return models.User{}, err
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.GetUserByID, "Получил пользователя по ID")
	if err != nil {
		return models.User{}, fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.GetUserByID,
		UserID: rand.Int(),
	})

	return user, nil
}

func (s *service) Create(ctx context.Context, user models.User) error {
	err := user.Validate()
	if err != nil {
		return errors.ErrFromValidate
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		return errors.ErrNotFound
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.Create, "Создал пользователя")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Create,
		UserID: rand.Int(),
	})

	return nil

}

func (s *service) Update(ctx context.Context, id int, updatedName string) error {
	err := utils.ValidateID(id)
	if err != nil {
		return errors.ErrFromValidateID
	}

	err = utils.ValidateStrEmpty(updatedName)
	if err != nil {
		return errors.ErrFromValidate
	}

	err = s.repository.Update(ctx, id, updatedName)
	if err != nil {
		return errors.ErrNotFound
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.Update, "Изменил пользователя")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Update,
		UserID: rand.Int(),
	})

	return nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := utils.ValidateID(id)
	if err != nil {
		return errors.ErrFromValidateID
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return errors.ErrNotFound
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.Delete, "Удалил пользователя")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.Delete,
		UserID: rand.Int(),
	})

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, id int, pass models.Auth) error {
	err := utils.ValidateID(id)
	if err != nil {
		return err
	}

	err = pass.ValidateAuthPassword()
	if err != nil {
		return err
	}

	user, err := s.repository.GetUser(ctx, id, models.User{
		OldPassword: pass.OldPassword,
	})
	if err != nil {
		return err
	}

	if strings.TrimSpace(user.Password) != strings.TrimSpace(pass.OldPassword) {
		log.Println("user.Password", user.Password)
		log.Println("pass.OldPassword", pass.OldPassword)
		return errors.ErrWrongPassword
	}

	err = s.repository.UpdatePassword(ctx, models.User{
		ID:          id,
		NewPassword: pass.NewPassword,
	})
	if err != nil {
		return err
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.UpdatePassword, "изменил пароль")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.UpdatePassword,
		UserID: rand.Int(),
	})
	return nil
}

func (s *service) UpdateAge(ctx context.Context, id int, req models.User) error {
	err := utils.ValidateID(id)
	if err != nil {
		return err
	}

	err = utils.ValidateInt(req.Age)
	if err != nil {
		return err
	}

	_, err = s.repository.GetUser(ctx, id, models.User{
		ID: id,
	})
	if err != nil {
		return fmt.Errorf("error from s.repository.SelectUser %w", err)
	}

	err = s.repository.UpdateAge(ctx, req.Age, id)
	if err != nil {
		return fmt.Errorf("error from s.repository.UpdateAge")
	}

	err = eventLogs.Audit(rand.Int(), eventLogs.UpdateAge, "изменил возраст")
	if err != nil {
		return fmt.Errorf("error from eventLogs.Audit %w", err)
	}

	s.eventsBus.Publish(eventBus.Event{
		Type:   eventLogs.UpdateAge,
		UserID: rand.Int(),
	})

	return nil
}
