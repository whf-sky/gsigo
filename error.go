package gsigo

import (
	"encoding/json"
	"fmt"
)

func ThrowError(message string, code ...int) error {
	errorCode := 0
	if len(code) > 0 {
		errorCode = code[0]
	}
	return &GsigoError{
		Code: errorCode,
		Message: message,
	}
}

type GsigoError struct {
	Code int
	Message string
	file string
	line int
}

func (ge GsigoError) Error() string {
	return fmt.Sprintf("%v:%v:%v: %v", ge.file, ge.line, ge.Code, ge.Message )
}

func (ge *GsigoError) getMessage () string {
	return ge.Message
}

func (ge *GsigoError) getCode () int {
	return ge.Code
}

func (ge *GsigoError) getLine () int {
	return ge.line
}

func (ge *GsigoError) getFile () string {
	return ge.file
}

func (ge *GsigoError) toJson () string {
	jsonByte,_ := json.Marshal(ge)
	return string(string(jsonByte))
}
