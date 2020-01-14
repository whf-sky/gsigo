package gsigo

import (
	"encoding/json"
	"fmt"
	"runtime"
)

//ThrowError throw error
func ThrowError(message string, code ...int) error {
	errorCode := 0
	if len(code) > 0 {
		errorCode = code[0]
	}
	_, file, line, _ := runtime.Caller(1)
	return &Error{
		Code: errorCode,
		Message: message,
		File:file,
		Line:line,
	}
}

type Error struct {
	Code int
	Message string
	File string
	Line int
}

func (ge Error) Error() string {
	return fmt.Sprintf("%v:%v:%v: %v", ge.File, ge.Line, ge.Code, ge.Message )
}

//getMessage get error message
func (ge *Error) getMessage () string {
	return ge.Message
}

//getCode get error code
func (ge *Error) getCode () int {
	return ge.Code
}

//getLine get error file line
func (ge *Error) getLine () int {
	return ge.Line
}

//getFile get error file
func (ge *Error) getFile () string {
	return ge.File
}

//toJson Error struct to json
func (ge *Error) toJson () string {
	jsonByte,_ := json.Marshal(ge)
	return string(string(jsonByte))
}
