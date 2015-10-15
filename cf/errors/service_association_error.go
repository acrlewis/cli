package errors

import "github.com/cloudfoundry/cli/cf/i18n"

type ServiceAssociationError struct {
}

func NewServiceAssociationError() error {
	return &ServiceAssociationError{}
}

func (err *ServiceAssociationError) Error() string {
	return i18n.T("Cannot delete service instance, service keys and bindings must first be deleted")
}
