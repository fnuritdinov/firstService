package models

import (
	_ "errors"
	"strconv"

	"github.com/fnuritdinov/firstService/pkg/errors"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (u *User) ValidateID() error {
	if u.ID < 1 {
		return errors.ErrorFromValidateID
	}
	return nil
}

func (u *User) ValidateStrEmpty() error {
	if len(u.Name) == 0 {
		return errors.ErrorFromValidateStrEmpty
	}
	return nil
}

func StrToInt(str string) int {
	n, _ := strconv.Atoi(str)
	return n
}

func IntToStr(n int) string {
	return strconv.Itoa(n)
}
