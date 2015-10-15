package errors

import (
	"fmt"

	"github.com/cloudfoundry/cli/cf/i18n"
)

type AsyncTimeoutError struct {
	url string
}

func NewAsyncTimeoutError(url string) error {
	return &AsyncTimeoutError{url: url}
}

func (err *AsyncTimeoutError) Error() string {
	return fmt.Sprintf(i18n.T("Error: timed out waiting for async job '{{.ErrURL}}' to finish",
		map[string]interface{}{"ErrURL": err.url}))
}
