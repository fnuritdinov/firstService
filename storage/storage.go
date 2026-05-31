package storage

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/fnuritdinov/firstService/models"
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

func (s *UserStorage) GetAll() ([]models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var users []models.User

	slFile, err := os.ReadFile(s.FileName)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(slFile, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) GetByID(id int) (models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	file, err := os.ReadFile(s.FileName)
	if err != nil {
		return models.User{}, errors.ErrorFromFile
	}

	var users []models.User

	err = json.Unmarshal(file, &users)
	if err != nil {
		return models.User{}, errors.ErrorParsingData
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, errors.ErrorNotFound
}

func (s *UserStorage) Create(user models.User) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var users []models.User
	file, err := os.ReadFile(s.FileName)
	if err != nil {
		return errors.ErrorFromFile
	}
	_ = json.Unmarshal(file, &users)

	users = append(users, user)

	newData, err := json.Marshal(users)
	if err != nil {
		return errors.ErrorParsingData
	}

	err = os.WriteFile(s.FileName, newData, 0644)
	if err != nil {
		return errors.ErrorFromFile
	}
	return nil
}

func (s *UserStorage) Update(id int, updatedUser string) error {

	s.Mu.Lock()
	defer s.Mu.Unlock()

	btData, err := os.ReadFile(s.FileName)
	if err != nil {
		return err
	}

	var users []models.User

	err = json.Unmarshal(btData, &users)
	if err != nil {
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
		return errors.ErrorParsingData
	}

	err = os.WriteFile(s.FileName, newData, 0644)
	if err != nil {
		return errors.ErrorFromFile
	}
	return nil
}

func (s *UserStorage) Delete(id int) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	bt, err := os.ReadFile(s.FileName)
	if err != nil {
		return errors.ErrorFromFile
	}

	var users []models.User

	err = json.Unmarshal(bt, &users)
	if err != nil {
		return errors.ErrorParsingData
	}

	for idx, value := range users {
		if value.ID == id {
			users = append(users[:idx], users[idx+1:]...)
			break
		}
	}
	btData, err := json.Marshal(users)
	if err != nil {
		return errors.ErrorParsingData
	}

	err = os.WriteFile(s.FileName, btData, 0644)
	if err != nil {
		return errors.ErrorFromFile
	}
	return nil
}
