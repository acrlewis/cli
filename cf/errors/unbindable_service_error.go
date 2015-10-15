package errors

import "github.com/cloudfoundry/cli/cf/i18n"

type UnbindableServiceError struct {
}

func NewUnbindableServiceError() error {
	return &UnbindableServiceError{}
}

func (err *UnbindableServiceError) Error() string {
	return i18n.T("This service doesn't support creation of keys.")
}
