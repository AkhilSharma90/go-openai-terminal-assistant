package system

type OperatingSystem int

// These are the constants representing different operating systems.
// The iota keyword represents successive untyped integer constants.
const (
	UnknownOperatingSystem OperatingSystem = iota
	LinuxOperatingSystem
	MacOperatingSystem
	WindowsOperatingSystem
)

// String is a method on the OperatingSystem type that returns a string representation of the operating system.
func (o OperatingSystem) String() string {
	switch o {
	case LinuxOperatingSystem:
		return "linux"
	case MacOperatingSystem:
		return "macOS"
	case WindowsOperatingSystem:
		return "windows"
	default:
		return "unknown"
	}
}
