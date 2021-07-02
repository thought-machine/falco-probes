package operatingsystem

// KernelPackage abstracts the implementation of resolving the required inputs for build a Falco Driver for a Kernel Package.
// The outputs of these are not guaranteed to be unique,
type KernelPackage interface {
	// GetKernelRelease returns the value to mock as the output of `uname -r`.
	GetKernelRelease() string
	// GetKernelVersion returns the value to mock as the output of `uname -v`.
	GetKernelVersion() string
	// GetKernelMachine returns the value to mock as the output of `uname -m`.
	GetKernelMachine() string
	// GetOSRelease returns the file contents to use as the mock of `/etc/os-release`.
	GetOSRelease() FileContents
	// GetKernelConfiguration returns the volume to mount as `/host/lib/modules/`.
	GetKernelConfiguration() Volume
	// GetKernelSources returns the volume to mount as `/usr/src/`.
	GetKernelSources() Volume
}

// FileContents represents the contents of a file.
type FileContents string

// Volume represents a reference to a Volume (structured collection of files).
type Volume string
