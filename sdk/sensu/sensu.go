package sensu

import (
	"fmt"
	"os"
)

// ErrorCode type
type ErrorCode int

// Error codes
const (
	GeneralGolangError ErrorCode = 129 // internal script error
	ConfigError        ErrorCode = 127 // unix config error, not enough parms, etc
	PermissionError    ErrorCode = 126 // not executable, etc
	RuntimeError       ErrorCode = 42  // self explantory
	DebugError         ErrorCode = 37  // You had the Alliance on you, criminals and savages… half the people on this

	// ship have been shot or wounded, including yourself, and you’re harboring known fugitives.
	UnknownError  ErrorCode = 3 // Would it save you a lot of time if I just gave up and went mad now?
	CriticalError ErrorCode = 2 // “The ships hung in the sky in much the same way that bricks don't.”
	WarningError  ErrorCode = 1 // this kinda sucks but don't get out of bed to deal with it
	OkError       ErrorCode = 0 // “We’re still flying”
)

// Error wraps a sensu error code and message
type Error struct {
	Code    ErrorCode
	Message string
}

// NewError creates a sensu error
func NewError(code ErrorCode, format string, a ...interface{}) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

func (e *Error) Error() string {
	return e.Message
}

// Exit application with a message
func Exit(code ErrorCode, a ...interface{}) {
	if len(a) > 0 {
		if format, ok := a[0].(string); ok && len(a) > 1 {
			fmt.Printf(format+"\n", a[1:]...)
		} else {
			fmt.Println(a...)
		}
	}
	os.Exit(int(code))
}

// ExitWithError will exit the application if err is not nil printing the error message
func ExitWithError(err error) {
	if err != nil {
		if e, ok := err.(*Error); ok {
			Exit(e.Code, e.Message)
		}
		Exit(RuntimeError, err.Error())
	}
}
