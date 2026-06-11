package models

import (
	_ "errors"

	"github.com/fnuritdinov/firstService/pkg/errors"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	IsActive bool   `json:"IsActive"`
}

func (u *User) Validate() error {
	if u.ID < 1 {
		return errors.ErrFromValidate
	}
	if len(u.Name) == 0 {
		return errors.ErrFromValidate
	}
	if u.Age < 1 {
		return errors.ErrFromValidate
	}
	if len(u.Password) == 0 {
		return errors.ErrFromValidate
	}
	return nil
}
