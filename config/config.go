package config

//COMPLETE

import (
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/akhilsharma90/terminal-assistant/system"
	"github.com/spf13/viper"
)

//config contains the configs for user and ai (defined in ai.go and user.go files resp.)
//the last one in the list is the analysis struct from the system package

type Config struct {
	ai     AiConfig         // ai config
	user   UserConfig       // user config
	system *system.Analysis // system config
}

// GetAiConfig returns the ai config that we've defined in the ai.go file
func (c *Config) GetAiConfig() AiConfig {
	return c.ai
}

// GetUserConfig returns the user config defined in the user.go file
func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

// GetSystemConfig returns the system config basically analysis struct defined in system package
func (c *Config) GetSystemConfig() *system.Analysis {
	return c.system
}

// NewConfig creates a new Config instance by reading the configuration from the file.
// It sets the default values for AI and user configurations if they are not present in the file.
func NewConfig() (*Config, error) {

	//we will set the values for the ai and the user config below but the system config has
	//quite a few values that we need to set and that's done with the help of the Analyse function
	//the analyse function is in the analyzer.go file in the system package

	system := system.Analyse()

	// Set the configuration file name, which is a field in the analysis struct in analyzer.go file
	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	//the analysis struct has a field for config file path, that's what we're updating here
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", system.GetHomeDirectory()))

	// in the 2 lines before we have set the values for the viper tool
	//now we can call the readinconfig function because we have the name of the config file
	//and the path so this function now actually reads the file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	

	// To be able to set the config, we need to set values for the ai field in the struct,
	//the user field and the system field and we set the values for all of them here
	return &Config{
		ai: AiConfig{
			key:         viper.GetString(openai_key),
			model:       viper.GetString(openai_model),
			proxy:       viper.GetString(openai_proxy),
			temperature: viper.GetFloat64(openai_temperature),
			maxTokens:   viper.GetInt(openai_max_tokens),
		},
		user: UserConfig{
			defaultPromptMode: viper.GetString(user_default_prompt_mode),
			preferences:       viper.GetString(user_preferences),
		},
		system: system,
	}, nil
}

//this is the writeConfig function of the config package that gets called in the ui.go file
// WriteConfig writes the configuration to the file and also returns a new Config instance.
func WriteConfig(key string, write bool) (*Config, error) {
	//we first call the analyse function from the system package to get the values for the
	//system config. there are 2 other configs apart from system - ai and user and that's what
	//we'll tackle next
	system := system.Analyse()

	// Set the AI default values for all the fields in the ai config struct
	viper.Set(openai_key, key)
	viper.Set(openai_model, openai.GPT3Dot5Turbo)
	viper.SetDefault(openai_proxy, "")
	viper.SetDefault(openai_temperature, 0.2)
	viper.SetDefault(openai_max_tokens, 1000)

	// Setting the user defaults by setting values for the fields in the userConfig struct
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	if write {
		//now that all the config values for system (with the help of the analyse func),
		//user config and ai config values are set, we can write them to the config file
		// Write the configuration to the file
		err := viper.SafeWriteConfigAs(system.GetConfigFile())
		if err != nil {
			return nil, err
		}
	}

	// There's a function just 2 funcs before this which is the NewConfig function
	//we're basically calling that function from here
	return NewConfig()
}
