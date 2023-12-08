package config

import (
	"os"
	"strings"
	"testing"

	"github.com/akhilsharma90/terminal-assistant/system"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	// TestConfig is a test function that runs subtests for NewConfig and WriteConfig.
	t.Run("NewConfig", testNewConfig)
	t.Run("WriteConfig", testWriteConfig)
}

// setupViper initializes the Viper configuration for testing purposes.
func setupViper(t *testing.T) {
	t.Helper()
	system := system.Analyse()

	viper.SetConfigName(strings.ToLower(system.GetApplicationName()))
	viper.AddConfigPath("/tmp/")
	viper.Set(openai_key, "test_key")
	viper.Set(openai_model, openai.GPT3Dot5Turbo)
	viper.Set(openai_proxy, "test_proxy")
	viper.Set(openai_temperature, 0.2)
	viper.Set(openai_max_tokens, 2000)
	viper.Set(user_default_prompt_mode, "exec")
	viper.Set(user_preferences, "test_preferences")

	require.NoError(t, viper.SafeWriteConfigAs("/tmp/terminal-assistant.json"))
}

// cleanup is a helper function used in testing to remove the "/tmp/terminal-assistant.json" file.
func cleanup(t *testing.T) {
	t.Helper()
	require.NoError(t, os.Remove("/tmp/terminal-assistant.json"))
}

// testNewConfig is a unit test function that tests the NewConfig function.
// It sets up the necessary configuration for testing, creates a new config,
// and asserts that the values of the config match the expected values.
func testNewConfig(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := NewConfig()
	require.NoError(t, err)

	assert.Equal(t, "test_key", cfg.GetAiConfig().GetKey())
	assert.Equal(t, openai.GPT3Dot5Turbo, cfg.GetAiConfig().GetModel())
	assert.Equal(t, "test_proxy", cfg.GetAiConfig().GetProxy())
	assert.Equal(t, 0.2, cfg.GetAiConfig().GetTemperature())
	assert.Equal(t, 2000, cfg.GetAiConfig().GetMaxTokens())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())
}

// testWriteConfig is a unit test function that tests the behavior of the WriteConfig function.
// It sets up the necessary configuration, writes a new configuration with test values,
// and then asserts that the written configuration matches the expected values.
func testWriteConfig(t *testing.T) {
	setupViper(t)
	defer cleanup(t)

	cfg, err := WriteConfig("new_test_key", false)
	require.NoError(t, err)

	assert.Equal(t, "new_test_key", cfg.GetAiConfig().GetKey())
	assert.Equal(t, openai.GPT3Dot5Turbo, cfg.GetAiConfig().GetModel())
	assert.Equal(t, "test_proxy", cfg.GetAiConfig().GetProxy())
	assert.Equal(t, 0.2, cfg.GetAiConfig().GetTemperature())
	assert.Equal(t, 2000, cfg.GetAiConfig().GetMaxTokens())
	assert.Equal(t, "exec", cfg.GetUserConfig().GetDefaultPromptMode())
	assert.Equal(t, "test_preferences", cfg.GetUserConfig().GetPreferences())

	assert.NotNil(t, cfg.GetSystemConfig())

	assert.Equal(t, "new_test_key", viper.GetString(openai_key))
	assert.Equal(t, openai.GPT3Dot5Turbo, viper.GetString(openai_model))
	assert.Equal(t, "test_proxy", viper.GetString(openai_proxy))
	assert.Equal(t, 0.2, viper.GetFloat64(openai_temperature))
	assert.Equal(t, 2000, viper.GetInt(openai_max_tokens))
	assert.Equal(t, "exec", viper.GetString(user_default_prompt_mode))
	assert.Equal(t, "test_preferences", viper.GetString(user_preferences))
}
