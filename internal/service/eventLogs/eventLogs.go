package eventLogs

import (
	"errors"
	"fmt"
	"os"
	"time"

	errs "github.com/fnuritdinov/firstService/pkg/errors"
)

type EventLogs struct {
	Name        string
	Method      string
	Description string
	CreatedAt   time.Time
}

const Create = "Сreate"
const GetAll = "GetAll"
const GetUserByID = "GetUserByID"
const Update = "Update"
const Delete = "Delete"

func Audit(name string, description string, method string) error {
	if name == "" {
		return errors.New("name is empty")
	}

	switch method {
	case Create:
		fmt.Printf("%s создал пользователя ", name)
	case GetAll:
		fmt.Printf("%s получил всех пользователей", name)
	case GetUserByID:
		fmt.Printf("%s получил пользователя по айди", name)
	case Update:
		fmt.Printf("%s изменил пользователя", name)
	case Delete:
		fmt.Printf("%s удалил пользователя", name)
	default:
		fmt.Println("unknown audit method")
	}

	log := EventLogs{
		Name:        name,
		Method:      method,
		Description: description,
		CreatedAt:   time.Now(),
	}

	file, err := os.OpenFile("data/logs.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return errs.ErrorFromFile
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("log audit %+v\n", log))
	if err != nil {
		return errs.ErrorFromFile
	}
	return nil
}
