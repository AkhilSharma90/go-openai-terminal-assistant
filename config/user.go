package config

//COMPLETE
// Constants for the user configuration keys.
const (
	user_default_prompt_mode = "USER_DEFAULT_PROMPT_MODE"
	user_preferences         = "USER_PREFERENCES"
)

//similar pattern followed here as ai.go file, we're using struct and constants
//to set and work with the values required, in this case it is user's configs

// UserConfig struct holds the user's configuration.
type UserConfig struct {
	// defaultPromptMode is the user's default prompt mode.
	defaultPromptMode string
	// preferences are the user's preferences.
	preferences string
}

//below are helper functions that help get those values

// GetDefaultPromptMode returns the user's default prompt mode.
func (c UserConfig) GetDefaultPromptMode() string {
	return c.defaultPromptMode
}

// GetPreferences returns the user's preferences.
func (c UserConfig) GetPreferences() string {
	return c.preferences
}
