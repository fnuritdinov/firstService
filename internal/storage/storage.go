package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/fnuritdinov/firstService/internal/models"
	"github.com/fnuritdinov/firstService/pkg/errors"
)

type UserStorage struct {
	mu       sync.Mutex
	fileName string
}

func NewUserStorage(filename string) *UserStorage {
	return &UserStorage{
		fileName: filename,
	}
}

func (s *UserStorage) Login(ctx context.Context, request models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var users []models.User

	byteDate, err := os.ReadFile(s.fileName)
	if err != nil {
		return errors.ErrFromFile
	}
	err = json.Unmarshal(byteDate, &users)
	if err != nil {
		return errors.ErrParsingData
	}

	for _, user := range users {
		if user.Name == request.Name {
			if user.Password != request.Password {
				return errors.ErrWrongPassword
			}
			return nil
		}
	}

	return errors.ErrNotFound
}
func (s *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []models.User

	slFile, err := os.ReadFile(s.fileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return nil, err
	}

	err = json.Unmarshal(slFile, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) Get(ctx context.Context) ([]models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []models.User

	slFile, err := os.ReadFile(s.fileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return nil, err
	}

	err = json.Unmarshal(slFile, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return nil, err
	}
	var activeUsers []models.User

	for _, user := range users {
		if user.IsActive {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers, nil
}

func (s *UserStorage) GetByID(ctx context.Context, id int) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.ReadFile(s.fileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return models.User{}, errors.ErrFromFile
	}

	var users []models.User

	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return models.User{}, errors.ErrParsingData
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, errors.ErrNotFound
}

func (s *UserStorage) Create(ctx context.Context, user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []models.User
	file, err := os.ReadFile(s.fileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return errors.ErrFromFile
	}
	_ = json.Unmarshal(file, &users)

	users = append(users, user)

	newData, err := json.Marshal(users)
	if err != nil {
		fmt.Errorf("error from json.Marshal %w", err)
		return errors.ErrParsingData
	}

	err = os.WriteFile(s.fileName, newData, 0644)
	if err != nil {
		fmt.Errorf("error from os.WriteFile %w", err)
		return errors.ErrFromFile
	}
	return nil
}

func (s *UserStorage) Update(ctx context.Context, id int, updatedUser string) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	btData, err := os.ReadFile(s.fileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile", err)
		return err
	}

	var users []models.User

	err = json.Unmarshal(btData, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return err
	}

	for i, user := range users {
		if user.ID == id {
			users[i].Name = updatedUser
			break
		}
	}

	newData, err := json.Marshal(users)
	if err != nil {
		fmt.Errorf("error from json.Marshal %w", err)
		return errors.ErrParsingData
	}

	err = os.WriteFile(s.fileName, newData, 0644)
	if err != nil {
		fmt.Errorf("error from os.WriteFile %w", err)
		return errors.ErrFromFile
	}
	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bt, err := os.ReadFile(s.fileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return errors.ErrFromFile
	}

	var users []models.User

	err = json.Unmarshal(bt, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return errors.ErrParsingData
	}

	for _, value := range users {
		if value.ID == id {
			if !value.IsActive {
				return errors.ErrIsActiveFalse
			}
			value.IsActive = true
		}
	}

	btData, err := json.Marshal(users)
	if err != nil {
		fmt.Errorf("error from json.Marshal %w", err)
		return errors.ErrParsingData
	}

	err = os.WriteFile(s.fileName, btData, 0644)
	if err != nil {
		fmt.Errorf("error from os.WriteFile %w", err)
		return errors.ErrFromFile
	}
	return nil
}
