package models

import "github.com/fnuritdinov/firstService/pkg/errors"

type Auth struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
	UserID      int    `json:"userID"`
}

func (a *Auth) ValidateAuth() error {
	if len(a.Login) < 1 {
		return errors.ErrFromValidate
	}
	if len(a.Password) < 1 {
		return errors.ErrFromValidate
	}
	if a.UserID < 1 {
		return errors.ErrFromValidate
	}
	return nil
}

func (a *Auth) ValidateAuthPassword() error {
	if len(a.NewPassword) < 1 {
		return errors.ErrFromValidate
	}

	if len(a.OldPassword) < 1 {
		return errors.ErrFromValidate
	}

	return nil
}
