package eventLogs

import (
	"fmt"
	"os"
	"time"

	errs "github.com/fnuritdinov/firstService/pkg/errors"
)

type EventLogs struct {
	UserID      int
	Method      string
	Description string
	CreatedAt   time.Time
}

const Create = "Сreate"
const GetAll = "GetAll"
const GetUserByID = "GetUserByID"
const Update = "Update"
const Delete = "Delete"
const Get = "Get"
const Login = "Login"

func Audit(userID int, method, description string) error {

	switch method {
	case Create:
		fmt.Println("%d создал пользователя ", userID)
	case GetAll:
		fmt.Printf("%d получил всех пользователей", userID)
	case GetUserByID:
		fmt.Printf("%d получил пользователя по айди", userID)
	case Update:
		fmt.Printf("%d изменил пользователя", userID)
	case Delete:
		fmt.Printf("%d удалил пользователя", userID)
	case Get:
		fmt.Printf("%d получил активных пользователя", userID)
	case Login:
		fmt.Printf("%d пользователь залогинился", userID)

	default:
		fmt.Println("unknown audit method")
	}

	log := EventLogs{
		UserID:      userID,
		Method:      method,
		Description: description,
		CreatedAt:   time.Now(),
	}

	file, err := os.OpenFile("data/logs.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return errs.ErrFromFile
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("log audit %+v\n", log))
	if err != nil {
		return errs.ErrFromFile
	}
	return nil
}
