package operatingsystem

// OperatingSystem abstracts the implementation of determining which kernel packages are available and the retrieval of them for an Operating System.
type OperatingSystem interface {
	// GetKernelPackageNames returns a list of all available Kernel Package names.
	GetKernelPackageNames() ([]string, error)
	// GetKernelPackageByName returns a "hydrated" KernelPackage for the given Kernel Package name.
	// "hydrated" means that the values are retrieved, so this function should perform the fetching of Kernel Sources, etc. for a KernelPackage.
	GetKernelPackageByName(name string) (KernelPackage, error)
}
