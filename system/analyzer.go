package system

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/akhilsharma90/terminal-assistant/run"

	"github.com/mitchellh/go-homedir"
)

const APPLICATION_NAME = "terminal-assistant"

type Analysis struct {
	operatingSystem OperatingSystem // The operating system type.
	distribution    string          // The specific distribution of the OS.
	shell           string          // The shell being used.
	homeDirectory   string          // The home directory path.
	username        string          // The username of the current user.
	editor          string          // The default editor set.
	configFile      string          // The configuration file path.
}

// GetApplicationName is a method that returns the application name.
func (a *Analysis) GetApplicationName() string {
	return APPLICATION_NAME
}

// GetOperatingSystem is a method that returns the operating system type.
func (a *Analysis) GetOperatingSystem() OperatingSystem {
	return a.operatingSystem
}

// GetDistribution is a method that returns the distribution of the operating system.
func (a *Analysis) GetDistribution() string {
	return a.distribution
}

// GetShell is a method that returns the shell being used.
func (a *Analysis) GetShell() string {
	return a.shell
}

// GetHomeDirectory is a method that returns the home directory path.
func (a *Analysis) GetHomeDirectory() string {
	return a.homeDirectory
}

// GetUsername is a method that returns the username of the current user.
func (a *Analysis) GetUsername() string {
	return a.username
}

// GetEditor is a method that returns the default editor set.
func (a *Analysis) GetEditor() string {
	return a.editor
}

// GetConfigFile is a method that returns the configuration file path.
func (a *Analysis) GetConfigFile() string {
	return a.configFile
}

// Analyse is a function that returns an Analysis object.
func Analyse() *Analysis {
	return &Analysis{
		operatingSystem: GetOperatingSystem(),
		distribution:    GetDistribution(),
		shell:           GetShell(),
		homeDirectory:   GetHomeDirectory(),
		username:        GetUsername(),
		editor:          GetEditor(),
		configFile:      GetConfigFile(),
	}
}

// GetOperatingSystem is a function that determines the operating system where the application is running.
// It uses the GOOS variable from the runtime package, which contains the operating system target.
func GetOperatingSystem() OperatingSystem {
	switch runtime.GOOS {
	case "linux":
		return LinuxOperatingSystem
	case "darwin":
		return MacOperatingSystem
	case "windows":
		return WindowsOperatingSystem
	default:
		return UnknownOperatingSystem
	}
}

// GetDistribution is a function that determines the distribution of the Linux operating system.
// It runs the 'lsb_release -sd' command to get the distribution information.
func GetDistribution() string {
	dist, err := run.RunCommand("lsb_release", "-sd")
	if err != nil {
		return ""
	}

	return strings.Trim(strings.Trim(dist, "\n"), "\"")
}

// GetShell is a function that determines the shell being used in the operating system.
// It runs the 'echo $SHELL' command to get the shell information.
func GetShell() string {
	shell, err := run.RunCommand("echo", os.Getenv("SHELL"))
	if err != nil {
		return ""
	}

	// Split the shell information by the '/' character, trim the newline and double quote characters,
	// and return the last element of the split, which is the shell name.
	split := strings.Split(strings.Trim(strings.Trim(shell, "\n"), "\""), "/")

	return split[len(split)-1]
}

// GetHomeDirectory is a function that returns the home directory path.
func GetHomeDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}

// GetUsername is a function that returns the username of the current user.
func GetUsername() string {
	name, err := run.RunCommand("echo", os.Getenv("USER"))
	if err != nil {
		return ""
	}

	return strings.Trim(name, "\n")
}

// GetEditor is a function that returns the default editor set.
func GetEditor() string {
	name, err := run.RunCommand("echo", os.Getenv("EDITOR"))
	if err != nil {
		return "nano"
	}

	return strings.Trim(name, "\n")
}

// GetConfigFile is a function that returns the configuration file path.
func GetConfigFile() string {
	return fmt.Sprintf(
		"%s/.config/%s.json",
		GetHomeDirectory(),
		strings.ToLower(APPLICATION_NAME),
	)
}
