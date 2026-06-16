package utils

import (
	"strconv"

	"github.com/fnuritdinov/firstService/pkg/errors"
)

func ValidateID(id int) error {
	if id < 1 {
		return errors.ErrFromValidateID
	}
	return nil
}

func ValidateStrEmpty(name string) error {
	if len(name) == 0 {
		return errors.ErrFromValidate
	}
	return nil
}

func StrToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func ValidateInt(num int) error {
	if num < 1 {
		return errors.ErrFromValidate
	}
	return nil
}
