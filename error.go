package gsigo

import (
	"encoding/json"
	"fmt"
	"runtime"
)

//抛出错误信息
func ThrowError(message string, code ...int) *Error {
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

//错误信息转成字符串
func (ge *Error) Error() string {
	return fmt.Sprintf("%v:%v:%v: %v", ge.File, ge.Line, ge.Code, ge.Message )
}

//获取错误信息
func (ge *Error) getMessage () string {
	return ge.Message
}

//获取错误码
func (ge *Error) getCode () int {
	return ge.Code
}

//获取行号
func (ge *Error) getLine () int {
	return ge.Line
}

//获取文件
func (ge *Error) getFile () string {
	return ge.File
}

//错误信息转为json
func (ge *Error) toJson () string {
	jsonByte,_ := json.Marshal(ge)
	return string(jsonByte)
}
