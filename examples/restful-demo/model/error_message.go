package model

import "fmt"

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"detail"`
}

const (
	ErrorCodeParameter = iota
	ErrorCodeReadBody
	ErrorCodeUnmarshalJSON
)

var ErrorCodes = map[int]string{
	ErrorCodeParameter:     "Parameter error",
	ErrorCodeReadBody:      "Read body error",
	ErrorCodeUnmarshalJSON: "Unmarshal JSON error",
}

func NewErrorMessage(code int, err error) *ErrorMessage {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	return &ErrorMessage{code, ErrorCodes[code], errStr}
}

func NewCodeErrorMarkDownDoc(errorCodes map[int]string) string {
	doc := "**Error code reference table**\n"
	doc += "| code | message |\n"
	doc += "| ---  | --- |\n"
	for code, message := range errorCodes {
		doc += fmt.Sprintf("| %d | %s |\n", code, message)
	}
	return doc
}
