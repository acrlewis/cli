package formatters

import "github.com/cloudfoundry/cli/cf/i18n"

func Allowed(allowed bool) string {
	if allowed {
		return i18n.T("allowed")
	} else {
		return i18n.T("disallowed")
	}
}
