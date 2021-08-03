subinclude("//build/defs:operating_system")
subinclude("//amazonlinux2:amazonlinux2_kernel_packages")

for kernel_package_name in amazonlinux2_kernel_packages:
    kp=kernel_package(
        kernel_package_name = kernel_package_name,
        operating_system_rule = "//amazonlinux2:amazonlinux2",
    )

    # falco_drivers(
    #     name = kernel_package_name,
    #     kernel_package = kp,
    # )
