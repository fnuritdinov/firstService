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
	Mu       sync.Mutex
	FileName string
}

func NewUserStorage(filename string) *UserStorage {
	return &UserStorage{
		FileName: filename,
	}
}

func (s *UserStorage) Login(ctx context.Context, request models.User) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var users []models.User

	byteDate, err := os.ReadFile(s.FileName)
	if err != nil {
		return errors.ErrorFromFile
	}
	err = json.Unmarshal(byteDate, &users)
	if err != nil {
		return errors.ErrorParsingData
	}

	for _, user := range users {
		if user.Name == request.Name {
			if user.Password != request.Password {
				return errors.ErrWrongPassword
			}
			return nil
		}
	}

	return errors.ErrorNotFound
}
func (s *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var users []models.User

	slFile, err := os.ReadFile(s.FileName)
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
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var users []models.User

	slFile, err := os.ReadFile(s.FileName)
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
		if user.IsActive == true {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers, nil
}

func (s *UserStorage) GetByID(id int) (models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	file, err := os.ReadFile(s.FileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return models.User{}, errors.ErrorFromFile
	}

	var users []models.User

	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return models.User{}, errors.ErrorParsingData
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, errors.ErrorNotFound
}

func (s *UserStorage) Create(ctx context.Context, user models.User) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var users []models.User
	file, err := os.ReadFile(s.FileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return errors.ErrorFromFile
	}
	_ = json.Unmarshal(file, &users)

	users = append(users, user)

	newData, err := json.Marshal(users)
	if err != nil {
		fmt.Errorf("error from json.Marshal %w", err)
		return errors.ErrorParsingData
	}

	err = os.WriteFile(s.FileName, newData, 0644)
	if err != nil {
		fmt.Errorf("error from os.WriteFile %w", err)
		return errors.ErrorFromFile
	}
	return nil
}

func (s *UserStorage) Update(id int, updatedUser string) error {

	s.Mu.Lock()
	defer s.Mu.Unlock()

	btData, err := os.ReadFile(s.FileName)
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
		return errors.ErrorParsingData
	}

	err = os.WriteFile(s.FileName, newData, 0644)
	if err != nil {
		fmt.Errorf("error from os.WriteFile %w", err)
		return errors.ErrorFromFile
	}
	return nil
}

func (s *UserStorage) Delete(id int) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	bt, err := os.ReadFile(s.FileName)
	if err != nil {
		fmt.Errorf("error from os.ReadFile %w", err)
		return errors.ErrorFromFile
	}

	var users []models.User

	err = json.Unmarshal(bt, &users)
	if err != nil {
		fmt.Errorf("error from json.Unmarshal %w", err)
		return errors.ErrorParsingData
	}

	for _, value := range users {
		if value.ID == id {
			if value.IsActive == false {
				return errors.ErrIsActiveFalse
			}
			value.IsActive = true
		}
	}

	btData, err := json.Marshal(users)
	if err != nil {
		fmt.Errorf("error from json.Marshal %w", err)
		return errors.ErrorParsingData
	}

	err = os.WriteFile(s.FileName, btData, 0644)
	if err != nil {
		fmt.Errorf("error from os.WriteFile %w", err)
		return errors.ErrorFromFile
	}
	return nil
}
