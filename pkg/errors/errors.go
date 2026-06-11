package errors

import "errors"

var ErrNotFound = errors.New("not found")
var ErrParsingData = errors.New("error on parsing")
var ErrFromFile = errors.New("error on file")
var ErrJSONEncode = errors.New("error on encode")
var ErrFromConvert = errors.New("error from convert")
var ErrFromValidate = errors.New("error from validate")
var ErrFromValidateID = errors.New("error from id")
var ErrJSONDecoder = errors.New("error on decode")
var ErrWrongPassword = errors.New("wrong password")
var ErrIsActiveFalse = errors.New("isActive is false")
