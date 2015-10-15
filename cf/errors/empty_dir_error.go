package errors

import "github.com/cloudfoundry/cli/cf/i18n"

type EmptyDirError struct {
	dir string
}

func NewEmptyDirError(dir string) error {
	return &EmptyDirError{dir: dir}
}

func (err *EmptyDirError) Error() string {
	return err.dir + i18n.T(" is empty")
}
