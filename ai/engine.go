package ai

//COMPLETE

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

//creating the main engine here for the application, the engine has 2 modes - exec and chat
//the messages sent for the exec mode are stored in the execMessages field and the chat ones
//are stored in the chatMessages field. the config field here is the config struct from config package
//client is the open ai client from the go-openai package. running shows if the engine is running
//channel is for the output steam from the engine, it's a struct in the output.go file
type Engine struct {
	mode         EngineMode                     // The mode of the engine either ExecEngineMode or ChatEngineMode
	config       *config.Config                 // The configuration settings for the engine
	client       *openai.Client                 // The OpenAI API client
	execMessages []openai.ChatCompletionMessage // Messages for executing commands
	chatMessages []openai.ChatCompletionMessage // Messages for chat interactions
	channel      chan EngineChatStreamOutput    // The channel for sending chat stream output
	pipe         string                         // The pipe is the same as pipe in regular software engineering, turns the output from previous into input for the new
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
	//running is kept as false since we have just made the engine, but it isn't running yet
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

//Four helper functions below to set and get values for the engine

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
	//engine has a channel field that can take messages of type EngineChatSteamOuput
	//which is a struct defined in the output.go file
	// Send an EngineChatStreamOutput with the interrupt flag set to true
	e.channel <- EngineChatStreamOutput{
		content:    "[Interrupt]",
		last:       true,
		interrupt:  true,
		executable: false,
	}

	// Set the running flag to false to stop the engine, we're setting the field of the engine
	e.running = false

	return e
}

//we want to clear the execMessages and chatMessages fields in the engine and set them with
//empty values of the same type and this is what we do with the clear function
// Clear clears the Engine messages based on the current mode.
func (e *Engine) Clear() *Engine {
	//if the mode is exec mode, then make exex messages empty
	if e.mode == ExecEngineMode {
		e.execMessages = []openai.ChatCompletionMessage{} // Clear execMessages by creating a new empty slice
	} else {
		//else the mode would be chat mode, because there are only two modes, so we make that empty
		e.chatMessages = []openai.ChatCompletionMessage{} // Clear chatMessages
	}

	return e
}

// Reset resets the Engine.
//clear and reset are similar but different. in clear we empty the messages from the mode that 
//we're in but with reset, we're emptying everything in the engine, like all msgs, irrespective
func (e *Engine) Reset() *Engine {
	e.execMessages = []openai.ChatCompletionMessage{}
	e.chatMessages = []openai.ChatCompletionMessage{}

	return e
}

// ExecCompletion execute a completion request to the OpenAI API and process the response.
func (e *Engine) ExecCompletion(input string) (*EngineExecOutput, error) {
	ctx := context.Background()

	//when engine is processing, we set it to true
	// Set the running flag to true
	e.running = true

	// Append user message to the chat messages
	//this method is defined a few lines below in this file only
	e.appendUserMessage(input)

	// Create a chat completion request to the OpenAI API
	//we will capture the response in the resp variable
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

	// the content is in the choices, message and content inside resp variable
	//we will fetch it and store it in content variable for easier access to it
	content := resp.Choices[0].Message.Content

	//we're maintaining the messages from user and from assistant (our terminal assistant)
	//in the chat and referring to them as user message and assistant message and appending them
	// Append assistant message to the chat messages
	//this method is defined a few files below in this file only
	e.appendAssistantMessage(content)

	var output EngineExecOutput

	//we're expecting that the open ai will send us back the command that we have to execute
	//this command will be part of the content string that we have extracted from the response above
	//to get the command from the content string, we're doing a bit of regex match

	// Unmarshal the content into the EngineExecOutput struct
	err = json.Unmarshal([]byte(content), &output)
	if err != nil {
		re := regexp.MustCompile(`\{.*?\}`)
		match := re.FindString(content)
		if match != "" {
			// Unmarshal the match into the EngineExecOutput struct
			//we have defined the output variable which is of type EngineExecOutput struct
			//we wil unmarshal the value received in match into this object
			err = json.Unmarshal([]byte(match), &output)
			if err != nil {
				return nil, err
			}
		} else {
			//if we're unable to unmarshal the command from the content string, we attach
			//the entire content object itself to the explanation field in the execoutput
			// If unable to unmarshal, create a default EngineExecOutput
			output = EngineExecOutput{
				Command:     "",
				Explanation: content,
				Executable:  false,
			}
		}
	}
//we return the output object, which is of type execstreamoutput, which is what we're supposed to retun from this function
	return &output, nil
}

// ChatCompletion execute a completion request to the OpenAI API and process the response in real-time.
func (e *Engine) ChatStreamCompletion(input string) error {
	ctx := context.Background()

	// Set the running flag to true when the engine is processing the output
	e.running = true

	// Append user message to chat messages
	e.appendUserMessage(input)

	// Create a chat completion request to the OpenAI API, just like we did in the execCompletion func.
	req := openai.ChatCompletionRequest{
		Model:     e.config.GetAiConfig().GetModel(),
		MaxTokens: e.config.GetAiConfig().GetMaxTokens(),
		Messages:  e.prepareCompletionMessages(),
		Stream:    true,
	}

	// We now send the request object to the CreateChatCompletionStream function to create a stream
	//the stream is accessible to us via the stream variable now
	stream, err := e.client.CreateChatCompletionStream(ctx, req)
	//handle the error from above
	if err != nil {
		return err
	}
	//defer the closing of the stream to the end of this function's execution
	defer stream.Close()

	//define output as a variable of type string
	var output string

	for {
//since we are streaming the output and the engine has been set to running until engine is streaming the contents
//we will place a condition to check engine running or in other words, streaming of output
//we want to add each message received in the stream to the output variable
		if e.running {
			//receive the output from the stream, stream being the variable defined above
			//capture that in the response variable, we will extract the actual content from this later
			resp, err := stream.Recv()

			// Check if completion is finished (i.e if we've reached end of the stream)
			//if completion is finished, we would have to set engine running as false
			//so that we can break out of the loop above
			if errors.Is(err, io.EOF) {
				//by default exeecutable is false because we're working with chat completion right now
				//but there's a chance that we're working with commands even in chat completion mode, so it's
				//recommended to handle this particular situation
				executable := false
				// Check if the output is executable, i.e if the engine is in executable mode
				if e.mode == ExecEngineMode {
					if !strings.HasPrefix(output, noexec) && !strings.Contains(output, "\n") {
						//if it is, then set executable to true
						executable = true
					}
				}

// Send last output to channel saying that this is the last message because we are in the 
//if condition where we're checking if this is the end of the stream and we're setting the
//exectuable value of the streamOutput with the actual value that's in the executable variable
//it could be true or false based on what we set it above, doesn't matter, whatever it is, we're setting it here
				e.channel <- EngineChatStreamOutput{
					content:    "",
					last:       true,
					executable: executable,
				}
//since we're in the condition where the end of file has been reached, we will
//set the engine running field to false so we can break out of the outer loop
				e.running = false

//The last message from the stream would be the one from open ai, so we're setting it using the appendAssistantMessage func.
// Append assistant message to chat messages
				e.appendAssistantMessage(output)

				return nil
			}
//now we're out of the end of file loop, but still inside the running loop
			if err != nil {
//if there's an error, we set the engine running to false
				e.running = false
				return err
			}

// The response from openai is in resp variable but we need the actual string, that's in content inside delta, choices
			delta := resp.Choices[0].Delta.Content
//add the contents of the delta variable to output, which is a string
//we will keep extracting the contents from the response and adding it to the output variable
//because we're working with a stream, meaning we won't get all our response in one go
			output += delta

// Send output to channel, saying that this isn't the last message and setting the actual content from delta to the content variable
//setting last msg to false because the engine is still running and we're out of the end of file loop
//keep the output stream going and adding delta content to it
			e.channel <- EngineChatStreamOutput{
				content: delta,
				last:    false,
			}

			// time.Sleep(time.Microsecond * 100)
		} else {
//if the engine isn't running, we will close the stream and return nil
			stream.Close()
			return nil
		}
	}
}

// appendUserMessage appends a user message to the chat messages in the Engine.
//this function has been called multiple times in the functions above
func (e *Engine) appendUserMessage(content string) *Engine {
	if e.mode == ExecEngineMode {
		//if the execute mode is exec mode, then we add the user's message to exec messages object
		e.execMessages = append(e.execMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	} else {
		//else we will be in chat mode, since there are only 2 modes and we will add the user's message 
		//to the chat messages object
		e.chatMessages = append(e.chatMessages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		})
	}

	return e
}

// appendAssistantMessage appends an assistant message to the chat messages in the Engine.
//this function works exactly the same way as the func. above and has the same logic
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
//we call this function in the execCompletion funciton and the chatStreamCompletion function to help
//create the chat completion messages we send to Open AI API.
func (e *Engine) prepareCompletionMessages() []openai.ChatCompletionMessage {
	// Initialize a slice of chat completion messages based on the struct chatCompletionMessage
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: e.prepareSystemPrompt(), // Prepare the system prompt, calling this function which is defined below
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
//check engine mode and accordingly append messages
	if e.mode == ExecEngineMode {
		messages = append(messages, e.execMessages...) // Append the user and assistant messages for execution mode.
	} else {
		messages = append(messages, e.chatMessages...) // Append the user and assistant messages for chat mode.
	}

	return messages
}

// preparePipePrompt prepares the pipe prompt.
func (e *Engine) preparePipePrompt() string {
	//pipe is the same as the software engineering concept, it helps in continuing the conversation
	//or the stream by making the output from previous process into the input for the next process
	return fmt.Sprintf("I will work on the following input: %s", e.pipe)
}

// prepareSystemPrompt prepares the system prompt based on the current mode.
func (e *Engine) prepareSystemPrompt() string {
	var bodyPart string
	// If the mode is ExecEngineMode, prepare the system prompt for execution mode else prepare the system prompt for chat mode.
	if e.mode == ExecEngineMode {
		//exec part is a different function and chat part is a separate function, both have their own prompts
		//that are defined in the respective functions
		bodyPart = e.prepareSystemPromptExecPart()
	} else {
		bodyPart = e.prepareSystemPromptChatPart()
	}

	// Return the system prompt with the context part, calling the context function defined below
	//context here basically means telling the system about the operating system, distribution and such things
	return fmt.Sprintf("%s\n%s", bodyPart, e.prepareSystemPromptContextPart())
}

// prepareSystemPromptExecPart prepares the system prompt for execution mode.
//preparing chat gpt to operate in command execution mode by giving it a long prompt
//and setting the context, we want it to generate commands and run them
func (e *Engine) prepareSystemPromptExecPart() string {
	return "You are terminal-assistant, a powerful terminal assistant generating a JSON containing a command line for my input.\n" +
		"You will always reply using the following json structure: {\"cmd\":\"the command\", \"exp\": \"some explanation\", \"exec\": true}.\n" +
		"Your answer will always only contain the json structure, never add any advice or supplementary detail or information, even if I asked the same question before.\n" +
		"The field cmd will contain a single line command (don't use new lines, use separators like && and ; instead).\n" +
		"The field exp will contain an short explanation of the command if you managed to generate an executable command, otherwise it will contain the reason of your failure.\n" +
		"The field exec will contain true if you managed to generate an executable command, false otherwise." +
		"\n" +
		"Examples:\n" +
		"Me: list all files in my home dir\n" +
		"terminal-assistant: {\"cmd\":\"ls ~\", \"exp\": \"list all files in your home dir\", \"exec\\: true}\n" +
		"Me: list all pods of all docker images\n" +
		"terminal-assistant: {\"cmd\":\"docker image ls\", \"exp\": \"list of all images from your docker instance\", \"exec\": true}\n" +
		"Me: how are you ?\n" +
		"terminal-assistant: {\"cmd\":\"\", \"exp\": \"I'm good thanks but I cannot generate a command for this. Use the chat mode to discuss.\", \"exec\": false}"
}

// prepareSystemPromptChatPart prepares the system prompt for chat mode.
//preparing chat gpt for the chat mode to help as an assistant by replying to question
//in this mode, we want it to just reply and not generate or execute any commands
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
//values such as operating system, distribution, home directory and shell will be sent as prompt to 
//open ai so that the commands it creates for us will be easily executable and will be relevant to us
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
