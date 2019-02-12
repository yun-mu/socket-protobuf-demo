package constant

import "errors"

const (
	// error
	ErrorMsgParamWrong = "param wrong"
)

var (
	// error
	ErrorOutOfRange    = errors.New("out of range")
	ErrorIDFormatWrong = errors.New("id format is wrong")
	ErrorNotFound      = errors.New("not found")
	ErrorHasExist      = errors.New("has exist")
	ErrorNotExist      = errors.New("not exist")
	ErrorParamWrong    = errors.New("param is wrong")
	ErrorUnAuth        = errors.New("un auth")
)
