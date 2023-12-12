package config

// Constants for AI configuration keys
const (
	openai_key         = "OPENAI_KEY"         // Key for OpenAI API
	openai_model       = "OPENAI_MODEL"       // Model to use for OpenAI API
	openai_proxy       = "OPENAI_PROXY"       // Proxy to use for OpenAI API
	openai_temperature = "OPENAI_TEMPERATURE" // Temperature setting for OpenAI API
	openai_max_tokens  = "OPENAI_MAX_TOKENS"  // Maximum tokens to generate for OpenAI API
)

// AiConfig represents the configuration for the AI.
type AiConfig struct {
	key         string
	model       string
	proxy       string
	temperature float64
	maxTokens   int
}

// GetKey returns the key for OpenAI API.
func (c AiConfig) GetKey() string {
	return c.key
}

// GetModel returns the model to use for OpenAI API.
func (c AiConfig) GetModel() string {
	return c.model
}

// GetProxy returns the proxy to use for OpenAI API.
func (c AiConfig) GetProxy() string {
	return c.proxy
}

// GetTemperature returns the temperature setting for OpenAI API.
func (c AiConfig) GetTemperature() float64 {
	return c.temperature
}

// GetMaxTokens returns the maximum tokens to generate for OpenAI API.
func (c AiConfig) GetMaxTokens() int {
	return c.maxTokens
}
