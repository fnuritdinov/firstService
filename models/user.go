package models

import (
	"errors"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var ErrorNotFound = errors.New("not found")
var ErrorParsingData = errors.New("error on parsing")
var ErrorFromFile = errors.New("error on file")
var ErrorJSONEncode = errors.New("error on encode")
var ErrorFromConvert = errors.New("error from convert")
var ErrorFromValidateID = errors.New("id must be greater than 0")
var ErrorFromValidateStrEmpty = errors.New("string is empty")
var ErrorJSONDecoder = errors.New("error on decode")

func ValidateID(id int) error {
	if id < 1 {
		return ErrorFromValidateID
	}
	return nil
}

func ValidateStrEmpty(name string) error {
	if len(name) == 0 {
		return ErrorFromValidateStrEmpty
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
