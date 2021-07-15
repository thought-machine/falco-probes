package operatingsystem

// OperatingSystem abstracts the implementation of determining which kernel packages are available and the retrieval of them for an Operating System.
type OperatingSystem interface {
	// GetName returns a unique string name for the implementation of this interface.
	GetName() string
	// GetKernelPackageNames returns a list of all available Kernel Package names.
	GetKernelPackageNames() ([]string, error)
	// GetKernelPackageByName returns a "hydrated" KernelPackage for the given Kernel Package name.
	// "hydrated" means that the values are retrieved, so this function should perform the fetching of Kernel Sources, etc. for a KernelPackage
	// and is the only place to return errors for those processes.
	GetKernelPackageByName(name string) (*KernelPackage, error)
}
