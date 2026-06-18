package models

import (
	_ "errors"

	"github.com/fnuritdinov/firstService/pkg/errors"
)

type User struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	IsActive    bool   `json:"IsActive"`
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	FromUserID  int    `json:"fromUserID"`
	ToUserID    int    `json:"toUserID"`
}

func (u *User) Validate() error {

	if len(u.Name) == 0 {
		return errors.ErrFromValidate
	}
	if u.Age < 1 {
		return errors.ErrFromValidate
	}
	if !u.IsActive {
		return errors.ErrFromValidate
	}

	return nil
}

type UserStats struct {
	TotalUsers    int     `json:"totalUsers"`
	ActiveUsers   int     `json:"activeUsers"`
	InActiveUsers int     `json:"inActiveUsers"`
	AverageAge    float64 `json:"averageAge"`
}
