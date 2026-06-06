package errors

import "errors"

var ErrorNotFound = errors.New("not found")
var ErrorParsingData = errors.New("error on parsing")
var ErrorFromFile = errors.New("error on file")
var ErrorJSONEncode = errors.New("error on encode")
var ErrorFromConvert = errors.New("error from convert")
var ErrorFromValidateID = errors.New("id must be greater than 0")
var ErrorFromValidateStrEmpty = errors.New("string is empty")
var ErrorJSONDecoder = errors.New("error on decode")
var ErrWrongPassword = errors.New("wrong password")
var ErrIsActiveFalse = errors.New("isActive is false")
