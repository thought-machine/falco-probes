package repository

// Repository abstracts the implementation of repositories which mirror Falco probes.
type Repository interface {
	// PublishProbe "publishes" (uploads/mirrors) the given probePath to the repository using the given driverVersion and probeName to organise probes.
	PublishProbe(driverVersion string, probeName string, probePath string) error
	// IsAlreadyMirrored returns whether or not the given probeName is already mirrored in the repository for the given driverVersion.
	IsAlreadyMirrored(driverVersion string, probeName string) bool
}
