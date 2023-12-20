package system

//COMPLETE

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/akhilsharma90/terminal-assistant/run"

	"github.com/mitchellh/go-homedir"
)

const APPLICATION_NAME = "terminal-assistant"

//analysis struct has the values for the system config that we use in the config package
type Analysis struct {
	operatingSystem OperatingSystem // The operating system type.
	distribution    string          // The specific distribution of the OS.
	shell           string          // The shell being used.
	homeDirectory   string          // The home directory path.
	username        string          // The username of the current user.
	editor          string          // The default editor set.
	configFile      string          // The configuration file path.
}

//below are a bunch of helper functions that'll help us get the values for the analysis struct
// GetApplicationName is a method that returns the application name, it's not a field in the struct
//but it's a constant we have defined above
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

// Analyse is a function that returns an Analysis object by calling functions for each of 
//the values required for the fields in the struct
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

//above we had 'methods' to get values from the struct
//below are the list of 'functions' that help us set values for the fields of the struct 
//so the difference between the methods above and the functions below is that the functions help us set the values
//whereas the methods only help us return or get the paritcular value from the struct

// GetOperatingSystem is a function that determines the operating system where the application is running.
// It uses the GOOS variable from the runtime package, which contains the operating system target.
//returns unknown if it doesn't match our list of OS
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
//we run a linux command to get the particular distribution we are on
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
//like bash, zsh etc. because depending on that our programs will change
func GetShell() string {
	shell, err := run.RunCommand("echo", os.Getenv("SHELL"))
	if err != nil {
		return ""
	}

	//we get the raw value and need to process it a bit to get the actual value we can store
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
	//we use the nano editor because that's always available if there's no set editor in the env variables
	if err != nil {
		return "nano"
	}

	return strings.Trim(name, "\n")
}

// GetConfigFile is a function that returns the configuration file path.
//the config file that we make to run this project is made in the .config folder and 
//is called terminal-assistant.json and that's what we're saying below, that it'll have json format
//and the APPLICATION_NAME that's a constant defined on top of this file
func GetConfigFile() string {
	return fmt.Sprintf(
		"%s/.config/%s.json",
		GetHomeDirectory(),
		strings.ToLower(APPLICATION_NAME),
	)
}
