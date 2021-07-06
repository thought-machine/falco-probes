package operatingsystem

// KernelPackage represents the required inputs for build a Falco Driver for a Kernel Package.
type KernelPackage struct {
	// OperatingSystem is the name of the Operating System that this KernelPackage belongs to.
	OperatingSystem string
	// Name is the name of the KernelPackage from the Operating System's perspective.
	Name string
	// KernelRelease is the value to mock as the output of `uname -r`.
	KernelRelease string
	// KernelVersion is the value to mock as the output of `uname -v`.
	KernelVersion string
	// KernelMachine is the value to mock as the output of `uname -m`.
	KernelMachine string
	// OSRelease is the file contents to use as the mock of `/etc/os-release`.
	OSRelease FileContents
	// KernelConfiguration is the volume to mount as `/host/lib/modules/`.
	KernelConfiguration Volume
	// KernelSources is the volume to mount as `/usr/src/`.
	KernelSources Volume
}

// FileContents represents the contents of a file.
type FileContents string

// Volume represents a reference to a Volume (structured collection of files).
type Volume string
