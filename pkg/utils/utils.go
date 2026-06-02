package utils

import (
	"strconv"

	"github.com/fnuritdinov/firstService/pkg/errors"
)

func ValidateID(id int) error {
	if id < 1 {
		return errors.ErrorFromValidateID
	}
	return nil
}

func ValidateStrEmpty(name string) error {
	if len(name) == 0 {
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
