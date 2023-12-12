package config

import (
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"

	"github.com/akhilsharma90/terminal-assistant/system"
	"github.com/spf13/viper"
)

type Config struct {
	ai     AiConfig         // ai config
	user   UserConfig       // user config
	system *system.Analysis // system config
}

// GetAiConfig returns the ai config
func (c *Config) GetAiConfig() AiConfig {
	return c.ai
}

// GetUserConfig returns the user config
func (c *Config) GetUserConfig() UserConfig {
	return c.user
}

// GetSystemConfig returns the system config
func (c *Config) GetSystemConfig() *system.Analysis {
	return c.system
}

// NewConfig creates a new Config instance by reading the configuration from the file.
// It sets the default values for AI and user configurations if they are not present in the file.
func NewConfig() (*Config, error) {
	system := system.Analyse()

	// Set the configuration file name and path
	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/", system.GetHomeDirectory()))

	// Read the configuration from the file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Create a new Config instance with the read configuration values
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

// WriteConfig writes the configuration to the file and returns a new Config instance.
func WriteConfig(key string, write bool) (*Config, error) {
	system := system.Analyse()

	// Set the AI defaults
	viper.Set(openai_key, key)
	viper.Set(openai_model, openai.GPT3Dot5Turbo)
	viper.SetDefault(openai_proxy, "")
	viper.SetDefault(openai_temperature, 0.2)
	viper.SetDefault(openai_max_tokens, 1000)

	// Set the user defaults
	viper.SetDefault(user_default_prompt_mode, "exec")
	viper.SetDefault(user_preferences, "")

	if write {
		// Write the configuration to the file
		err := viper.SafeWriteConfigAs(system.GetConfigFile())
		if err != nil {
			return nil, err
		}
	}

	// Return a new Config instance
	return NewConfig()
}
