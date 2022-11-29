package mock

import (
	"github.com/thought-machine/falco-probes/pkg/operatingsystem/cos/buildid"
)

// BuildIDValidator implements ValidatorInterface.
type BuildIDValidator struct {
	buildid.ValidatorInterface
}

// FilterInvalid just returns all build IDs.
func (v BuildIDValidator) FilterInvalid(buildIDsIn []string) ([]string, error) {
	return buildIDsIn, nil
}
