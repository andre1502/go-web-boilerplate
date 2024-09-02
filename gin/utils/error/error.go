package error

import (
	"fmt"
	"runtime"
)

type ErrorC struct {
	Func  string
	Param Param
	Err   error
}

type Param struct {
	I18nMsg string
	Data    map[string]any
}

func (ec ErrorC) Error() string {
	return fmt.Sprintf("[%s][[%s] -> %v] - %v", ec.Func, ec.Param.I18nMsg, ec.Param.Data, ec.Err)
}

func Fail(funcName string, message string, data map[string]any, err error) ErrorC {
	return ErrorC{
		Func: funcName,
		Param: Param{
			I18nMsg: message,
			Data:    data,
		},
		Err: err,
	}
}

func FuncName() string {
	pc, _, line, _ := runtime.Caller(1)

	return fmt.Sprintf("[%s:%d]", runtime.FuncForPC(pc).Name(), line)
}

func ParseError(err error) (errorC ErrorC) {
	var ok bool

	if err == nil {
		return ErrorC{Err: err}
	}

	if errorC, ok = err.(ErrorC); !ok {
		panic("errorC must model.ErrorC type")
	}

	return errorC
}
