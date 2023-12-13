package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/akhilsharma90/terminal-assistant/config"
	"github.com/akhilsharma90/terminal-assistant/system"

	"github.com/sashabaranov/go-openai"
)

const noexec = "[noexec]"

//creating the main engine here for ai, the struct has fields for mode, config etc.

type Engine struct {
	mode         EngineMode                     // The mode of the engine either ExecEngineMode or ChatEngineMode
	config       *config.Config                 // The configuration settings for the engine
	client       *openai.Client                 // The OpenAI API client
	execMessages []openai.ChatCompletionMessage // Messages for executing commands
	chatMessages []openai.ChatCompletionMessage // Messages for chat interactions
	channel      chan EngineChatStreamOutput    // The channel for sending chat stream output
	pipe         string                         // The pipe for communication with the engine
	running      bool                           // Indicates whether the engine is running or not
}

// NewEngine creates a new instance of the Engine struct.
// It takes the mode (EngineMode) and config (*config.Config) as parameters.
//and returns an instance of type Engine struct
func NewEngine(mode EngineMode, config *config.Config) (*Engine, error) {
	var client *openai.Client

	// Check if a proxy is configured in the AI config
	if config.GetAiConfig().GetProxy() != "" {

		// Create a client configuration with the API key
		clientConfig := openai.DefaultConfig(config.GetAiConfig().GetKey())

		// Parse the proxy URL
		proxyUrl, err := url.Parse(config.GetAiConfig().GetProxy())
		if err != nil {
			return nil, err
		}

		// Create a transport with the proxy URL
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		// Set the transport in the client configuration
		clientConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

		// Create a new client with the configured client configuration
		client = openai.NewClientWithConfig(clientConfig)
	} else {
		// Create a new client with the API key
		client = openai.NewClient(config.GetAiConfig().GetKey())
	}

	// Create a new instance of the Engine struct with the provided parameters
	return &Engine{
		mode:         mode,
		config:       config,
		client:       client,
		execMessages: make([]openai.ChatCompletionMessage, 0),
		chatMessages: make([]openai.ChatCompletionMessage, 0),
		channel:      make(chan EngineChatStreamOutput),
		pipe:         "",
		running:      false,
	}, nil
}

//Four helper functions below to get values set for the engine

// SetMode sets the mode of the Engine.
func (e *Engine) SetMode(mode EngineMode) *Engine {
	e.mode = mode

	return e
}

// GetMode returns the mode of the Engine.
func (e *Engine) GetMode() EngineMode {
	return e.mode
}

// GetChannel returns the channel of the Engine.
func (e *Engine) GetChannel() chan EngineChatStreamOutput {
	return e.channel
}

// SetPipe sets the pipe of the Engine.
func (e *Engine) SetPipe(pipe string) *Engine {
	e.pipe = pipe

	return e
}

// Interrupt interrupts the Engine operation.
func (e *Engine) Interrupt() *Engine {
	// Send an EngineChatStreamOutput with the interrupt flag set to true
	e.channel <- EngineChatStreamOutput{
		content:    "[Interrupt]",
		last:       true,
		interrupt:  true,
		executable: false,
	}

	// Set the running flag to false
	e.running = false

	return e
}

// Clear clears the Engine messages based on the current mode.
func (e *Engine) Clear() *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = []openai.ChatCompletionMessage{} // Clear execMessages by creating a new empty slice
	} else {
		e.chatMessages = []openai.ChatCompletionMessage{} // Clear chatMessages
	}

	return e
}

// Reset resets the Engine.
func (e *Engine) Reset() *Engine {
	e.execMessages = []openai.ChatCompletionMessage{}
	e.chatMessages = []openai.ChatCompletionMessage{}

	return e
}

// ExecCompletion execute a completion request to the OpenAI API and process the response.
func (e *Engine) ExecCompletion(input string) (*EngineExecOutput, error) {
	ctx := context.Background()

	// Set the running flag to true
	e.running = true

	// Append user message to the chat messages
	e.appendUserMessage(input)

	// Create a chat completion request to the OpenAI API
	resp, err := e.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:     e.config.GetAiConfig().GetModel(),
			MaxTokens: e.config.GetAiConfig().GetMaxTokens(),
			Messages:  e.prepareCompletionMessages(),
		},
	)
	if err != nil {
		return nil, err
	}

	// Get the assistant message from the response
	content := resp.Choices[0].Message.Content

	// Append assistant message to the chat messages
	e.appendAssistantMessage(content)

	var output EngineExecOutput

	// Unmarshal the content into the EngineExecOutput struct
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		re := regexp.MustCompile(`\{.*?\}`)
		match := re.FindString(content)
		if match != "" {
			// Unmarshal the match into the EngineExecOutput struct
			err = json.Unmarshal([]byte(match), &output)
			if err != nil {
				return nil, err
			}
		} else {
			// If unable to unmarshal, create a default EngineExecOutput
			output = EngineExecOutput{
				Command:     "",
				Explanation: content,
				Executable:  false,
			}
		}
	}

	return &output, nil
}

// ChatCompletion execute a completion request to the OpenAI API and process the response in real-time.
func (e *Engine) ChatStreamCompletion(input string) error {
	ctx := context.Background()

	// Set the running flag to true
	e.running = true

	// Append user message to chat messages
	e.appendUserMessage(input)

	// Create a chat completion request to the OpenAI API
	req := openai.ChatCompletionRequest{
		Model:     e.config.GetAiConfig().GetModel(),
		MaxTokens: e.config.GetAiConfig().GetMaxTokens(),
		Messages:  e.prepareCompletionMessages(),
		Stream:    true,
	}

	// Create chat completion stream
	stream, err := e.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	defer stream.Close()

	var output string

	for {
		if e.running {
			resp, err := stream.Recv()

			// Check if completion is finished
			if errors.Is(err, io.EOF) {
				executable := false
				// Check if the output is executable
				if e.mode == ExecEngineMode {
					if !strings.HasPrefix(output, noexec) && !strings.Contains(output, "\n") {
						executable = true
					}
				}

				// Send last output to channel
				e.channel <- EngineChatStreamOutput{
					content:    "",
					last:       true,
					executable: executable,
				}
				e.running = false

				// Append assistant message to chat messages
				e.appendAssistantMessage(output)

				return nil
			}

			if err != nil {
				e.running = false
				return err
			}

			// Get assistant message from response
			delta := resp.Choices[0].Delta.Content

			output += delta

			// Send output to channel
			e.channel <- EngineChatStreamOutput{
				content: delta,
				last:    false,
			}

			// time.Sleep(time.Microsecond * 100)
		} else {
			stream.Close()
			return nil
		}
	}
}

// appendUserMessage appends a user message to the chat messages in the Engine.
func (e *Engine) appendUserMessage(content string) *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	} else {
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	}

	return e
}

// appendAssistantMessage appends an assistant message to the chat messages in the Engine.
func (e *Engine) appendAssistantMessage(content string) *Engine {
	if e.mode == ExecEngineMode {
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	} else {
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
	}

	return e
}

// prepareCompletionMessages prepares the chat completion messages to be sent to the OpenAI API.
func (e *Engine) prepareCompletionMessages() []openai.ChatCompletionMessage {
	// Create a slice of chat completion messages
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: e.prepareSystemPrompt(), // Prepare the system prompt.
		},
	}

	// If the pipe is not empty, append a pipe prompt to the messages.
	if e.pipe != "" {
		messages = append(
			messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: e.preparePipePrompt(), // Prepare the pipe prompt.
			},
		)
	}

	if e.mode == ExecEngineMode {
		messages = append(messages, e.execMessages...) // Append the user and assistant messages for execution mode.
	} else {
		messages = append(messages, e.chatMessages...) // Append the user and assistant messages for chat mode.
	}

	return messages
}

// preparePipePrompt prepares the pipe prompt.
func (e *Engine) preparePipePrompt() string {
	return fmt.Sprintf("I will work on the following input: %s", e.pipe)
}

// prepareSystemPrompt prepares the system prompt based on the current mode.
func (e *Engine) prepareSystemPrompt() string {
	var bodyPart string
	// If the mode is ExecEngineMode, prepare the system prompt for execution mode else prepare the system prompt for chat mode.
	if e.mode == ExecEngineMode {
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	// Return the system prompt with the context part.
	return fmt.Sprintf("%s\n%s", bodyPart, e.prepareSystemPromptContextPart())
}

// prepareSystemPromptExecPart prepares the system prompt for execution mode.
func (e *Engine) prepareSystemPromptExecPart() string {
	return "Your are terminal-assistant, a powerful terminal assistant generating a JSON containing a command line for my input.\n" +
		"You will always reply using the following json structure: {\"cmd\":\"the command\", \"exp\": \"some explanation\", \"exec\": true}.\n" +
		"Your answer will always only contain the json structure, never add any advice or supplementary detail or information, even if I asked the same question before.\n" +
		"The field cmd will contain a single line command (don't use new lines, use separators like && and ; instead).\n" +
		"The field exp will contain an short explanation of the command if you managed to generate an executable command, otherwise it will contain the reason of your failure.\n" +
		"The field exec will contain true if you managed to generate an executable command, false otherwise." +
		"\n" +
		"Examples:\n" +
		"Me: list all files in my home dir\n" +
		"terminal-assistant: {\"cmd\":\"ls ~\", \"exp\": \"list all files in your home dir\", \"exec\\: true}\n" +
		"Me: list all pods of all namespaces\n" +
		"terminal-assistant: {\"cmd\":\"kubectl get pods --all-namespaces\", \"exp\": \"list pods form all k8s namespaces\", \"exec\": true}\n" +
		"Me: how are you ?\n" +
		"terminal-assistant: {\"cmd\":\"\", \"exp\": \"I'm good thanks but I cannot generate a command for this. Use the chat mode to discuss.\", \"exec\": false}"
}

// prepareSystemPromptChatPart prepares the system prompt for chat mode.
func (e *Engine) prepareSystemPromptChatPart() string {
	return "You are a powerful terminal assistant.\n" +
		"You will answer in the most helpful possible way.\n" +
		"Always format your answer in markdown format.\n\n" +
		"For example:\n" +
		"Me: What is 2+2 ?\n" +
		"terminal-assistant: The answer for `2+2` is `4`\n" +
		"Me: +2 again ?\n" +
		"terminal-assistant: The answer is `6`\n"
}

// prepareSystemPromptContextPart prepares the system prompt context part.
func (e *Engine) prepareSystemPromptContextPart() string {
	part := "My context: "

	if e.config.GetSystemConfig().GetOperatingSystem() != system.UnknownOperatingSystem {
		// If the operating system is not unknown, append the operating system to the context part.
		part += fmt.Sprintf("my operating system is %s, ", e.config.GetSystemConfig().GetOperatingSystem().String())
	}
	if e.config.GetSystemConfig().GetDistribution() != "" {
		// If the distribution is not empty, append the distribution to the context part.
		part += fmt.Sprintf("my distribution is %s, ", e.config.GetSystemConfig().GetDistribution())
	}
	if e.config.GetSystemConfig().GetHomeDirectory() != "" {
		// If the home directory is not empty, append the home directory to the context part.
		part += fmt.Sprintf("my home directory is %s, ", e.config.GetSystemConfig().GetHomeDirectory())
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		// If the shell is not empty, append the shell to the context part.
		part += fmt.Sprintf("my shell is %s, ", e.config.GetSystemConfig().GetShell())
	}
	if e.config.GetSystemConfig().GetShell() != "" {
		// If the editor is not empty, append the editor to the context part.
		part += fmt.Sprintf("my editor is %s, ", e.config.GetSystemConfig().GetEditor())
	}
	part += "take this into account. "

	// If the preferences are not empty, append the preferences to the context part.
	if e.config.GetUserConfig().GetPreferences() != "" {
		part += fmt.Sprintf("Also, %s.", e.config.GetUserConfig().GetPreferences())
	}

	return part
}
