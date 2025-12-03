package commands

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ExitCodeBaseError                = fmt.Errorf("exit with code")
	UnexpectedErrorReturnedByCommand = fmt.Errorf("unexpected error returned by command")
)

const (
	OkExitCode           = 0
	DefaultErrorExitCode = 1
)

// ExitCodeError is the error that any command should return doesn't matter if the exit code is 0 or anything else.
// This way, we allow the main flow to execute logic and print error messages in a custom way to be backwards
// compatible.
type ExitCodeError struct {
	Cause    error
	ExitCode int
}

func (ece *ExitCodeError) Error() string {
	return fmt.Errorf("%w: %d", ExitCodeBaseError, ece.ExitCode).Error()
}

func ExtractExitCode(err error) (exitCode int, rootCause error) {
	if err == nil {
		return OkExitCode, nil
	}
	var expected *ExitCodeError
	if !errors.As(err, &expected) {
		return DefaultErrorExitCode, fmt.Errorf("%w: %s: %s", UnexpectedErrorReturnedByCommand, reflect.TypeOf(err), err)
	}
	return expected.ExitCode, expected.Cause
}
