package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/fnuritdinov/firstService/models"
)

type UserStorage struct {
	Mu       sync.Mutex
	FileName string
}

func (s *UserStorage) GetAll() ([]models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	var users []models.User

	slFile, err := os.ReadFile(s.FileName)
	if err != nil {
		return []models.User{}, err
	}

	err = json.Unmarshal(slFile, &users)
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (s *UserStorage) GetByID(id int) (*models.User, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	file, err := os.Open(s.FileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		field := strings.Fields(line)

		if len(field) < 2 {
			continue
		}

		num, err := strconv.Atoi(field[0])
		if err != nil {
			continue
		}
		if num == id {
			user := &models.User{
				ID:   num,
				Name: field[1],
			}
			return user, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, errors.New("user not found")

}

func (s *UserStorage) Create(user models.User) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	file, err := os.OpenFile("users.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d %s", user.ID, user.Name))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) Update(id int, updatedUser models.User) error {

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

	found := false

	for i, user := range users {
		if user.ID == id {
			users[i].Name = updatedUser.Name
			found = true
			break
		}
	}

	newData, err := json.Marshal(users)
	if err != nil {
		return err
	}

	if !found {
		return errors.New("user not found")
	}

	err = os.WriteFile(s.FileName, newData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStorage) Delete(id int) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	bt, err := os.ReadFile(s.FileName)
	if err != nil {
		return err
	}

	var users []models.User

	err = json.Unmarshal(bt, &users)
	if err != nil {
		return err
	}

	for idx, value := range users {
		if value.ID == id {
			users = append(users[:idx], users[idx+1:]...)
			break
		}
	}
	return nil
}
