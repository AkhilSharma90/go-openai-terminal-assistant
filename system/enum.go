package system

//COMPLETE
//file deals with multiple operating systems

type OperatingSystem int

// These are the constants representing different operating systems.
// The iota keyword represents successive untyped integer constants.
const (
	UnknownOperatingSystem OperatingSystem = iota
	LinuxOperatingSystem
	MacOperatingSystem
	WindowsOperatingSystem
)

//this particular method we call in the engine.go file helps to return the operating system in the string format
//because if you check, in the analyzer.go file, the GetOperatingSystem function returns the value, not string
//and if you notice, that's the value being used here in this function to determine which string to return
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
