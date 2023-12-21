package ui

import (
	"fmt"
	"strings"

	"github.com/akhilsharma90/terminal-assistant/ai"
	"github.com/akhilsharma90/terminal-assistant/config"
	"github.com/akhilsharma90/terminal-assistant/history"
	"github.com/akhilsharma90/terminal-assistant/run"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/spf13/viper"
)

// UiState is a struct that represents the state of the user interface.
type UiState struct {
	error       error      // Any error that occurred.
	runMode     RunMode    // The mode in which the program is running.
	promptMode  PromptMode // The mode of the prompt.
	configuring bool       // Whether the program is in configuration mode.
	querying    bool       // Whether the program is in querying mode.
	confirming  bool       // Whether the program is in confirming mode.
	executing   bool       // Whether the program is in executing mode.
	args        string     // The arguments passed to the program.
	pipe        string     // The pipe used by the program.
	buffer      string     // The buffer of the program.
	command     string     // The command being executed by the program.
}

// UiDimensions is a struct that represents the dimensions of the user interface.
type UiDimensions struct {
	width  int
	height int
}

// UiComponents is a struct that represents the components of the user interface.
//there are three components we will work with and they each have their separate files
type UiComponents struct {
	prompt   *Prompt   // The prompt of the user interface.
	renderer *Renderer // The renderer of the user interface.
	spinner  *Spinner  // The spinner of the user interface.
}

// Ui is a struct that represents the user interface.
type Ui struct {
	state      UiState          // The state is of type UiState, a struct we have defined above
	dimensions UiDimensions     // UiDimensions is the struct defined above
	components UiComponents     // This is of type UiComponents, a struct we have defined above
	config     *config.Config   // Config struct in the config package has system config, user config etc.
	engine     *ai.Engine       // Engine is actually a struct we have in the ai package of this project
	history    *history.History // History is a struct in the history package of this project
}

//this function gets called from main.go
// NewUi is a function that creates a new Ui instance.
func NewUi(input *UiInput) *Ui {
	// Create a new Ui instance with the input run mode and prompt mode, a new prompt, renderer, and spinner, and a new history.
	//config and engine are not yet initialized
	return &Ui{
		state: UiState{
			error:       nil,
			runMode:     input.GetRunMode(),
			promptMode:  input.GetPromptMode(),
			configuring: false,
			querying:    false,
			confirming:  false,
			executing:   false,
			args:        input.GetArgs(),
			pipe:        input.GetPipe(),
			buffer:      "",
			command:     "",
		},
		dimensions: UiDimensions{
			150,
			150,
		},
		components: UiComponents{
			prompt: NewPrompt(input.GetPromptMode()), //this function is defined in prompt.go
			renderer: NewRenderer( //render.go has this function
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(150),
			),
			spinner: NewSpinner(), //spinner.go has this func
		},
		history: history.NewHistory(), //calls the helper function NewHistory in the history package 
	}
}



// Init initializes the UI and returns a tea.Cmd that represents the initial command to be executed.
// It loads the configuration, handles any errors, and determines whether to start in REPL mode or CLI mode.
func (u *Ui) Init() tea.Cmd {
	// Load the configuration
	config, err := config.NewConfig()
	if err != nil {
		// Handle the case when the configuration file is not found
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//if the user hasn't provided the config file, we are setting our terminal into
			//repl mode and accepting values
			if u.state.runMode == ReplMode {
				// If running in REPL mode, sequence the commands to clear the screen and start the configuration
				return tea.Sequence(
					tea.ClearScreen,
					u.startConfig(),
				)
			} else {
				// If running in CLI mode, start the configuration
				return u.startConfig()
			}
		} else {
			// Handle other errors by printing the error message and quitting the program
			//if user doesn't get config file and is also not able to accept config
			//in repl or cli mode, we send error
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderError(err.Error())),
				tea.Quit,
			)
		}
	}

	// Determine whether to start in REPL mode or CLI mode
	if u.state.runMode == ReplMode {
		// Start in REPL mode
		return u.startRepl(config)
	} else {
		// Start in CLI mode
		return u.startCli(config)
	}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Update is a method of the Ui struct that handles updating the UI based on the received message.
func (u *Ui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds       []tea.Cmd
		promptCmd  tea.Cmd
		spinnerCmd tea.Cmd
	)

	switch msg := msg.(type) {
	// Handle spinner tick message
	case spinner.TickMsg:
		if u.state.querying {
			u.components.spinner, spinnerCmd = u.components.spinner.Update(msg)
			cmds = append(
				cmds,
				spinnerCmd,
			)
		}
	// Handle window size message
	case tea.WindowSizeMsg:
		u.dimensions.width = msg.Width
		u.dimensions.height = msg.Height
		u.components.renderer = NewRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(u.dimensions.width),
		)
	// Handle keyboard input
	case tea.KeyMsg:
		switch msg.Type {
	//we are fixing the actions based on the key pressed by the user, like ctrlc, tab etc.
		// Quit the program
		case tea.KeyCtrlC:
			return u, tea.Quit
		// Navigate command history with up and down keys
		case tea.KeyUp, tea.KeyDown:
			if !u.state.querying && !u.state.confirming {
				var input *string
				if msg.Type == tea.KeyUp {
					input = u.history.GetPrevious()
				} else {
					input = u.history.GetNext()
				}
				if input != nil {
					u.components.prompt.SetValue(*input)
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)
					cmds = append(
						cmds,
						promptCmd,
					)
				}
			}
		// Switch between chat and execution mode
		case tea.KeyTab:
			if !u.state.querying && !u.state.confirming {
				if u.state.promptMode == ChatPromptMode {
					u.state.promptMode = ExecPromptMode
					u.components.prompt.SetMode(ExecPromptMode)
					u.engine.SetMode(ai.ExecEngineMode)
				} else {
					u.state.promptMode = ChatPromptMode
					u.components.prompt.SetMode(ChatPromptMode)
					u.engine.SetMode(ai.ChatEngineMode)
				}
				u.engine.Reset()
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		// [User presses enter] Process user input
		case tea.KeyEnter:
			if u.state.configuring {
				return u, u.finishConfig(u.components.prompt.GetValue())
			}
			if !u.state.querying && !u.state.confirming {
				input := u.components.prompt.GetValue()
				if input != "" {
					inputPrint := u.components.prompt.AsString()
					u.history.Add(input)
					u.components.prompt.SetValue("")
					u.components.prompt.Blur()
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)
					if u.state.promptMode == ChatPromptMode {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							u.startChatStream(input),
							u.awaitChatStream(),
						)
					} else {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(inputPrint),
							u.startExec(input),
							u.components.spinner.Tick,
						)
					}
				}
			}
		// [user presses ctrl + h] Show help message
		case tea.KeyCtrlH:
			if !u.state.configuring && !u.state.querying && !u.state.confirming {
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
					textinput.Blink,
				)
			}
		// [Ctrl + l] Clear the screen
		case tea.KeyCtrlL:
			if !u.state.querying && !u.state.confirming {
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}
		// [ctrl + R] Reset the program
		case tea.KeyCtrlR:
			if !u.state.querying && !u.state.confirming {
				u.history.Reset()
				u.engine.Reset()
				u.components.prompt.SetValue("")
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					tea.ClearScreen,
					textinput.Blink,
				)
			}
		// [Ctrl + S] Edit settings
		case tea.KeyCtrlS:
			if !u.state.querying && !u.state.confirming && !u.state.configuring && !u.state.executing {
				u.state.executing = true
				u.state.buffer = ""
				u.state.command = ""
				u.components.prompt.Blur()
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					u.editSettings(), //calling the editSettings function defined below
				)
			}
		default:
			if u.state.confirming {
				if strings.ToLower(msg.String()) == "y" {
					u.state.confirming = false
					u.state.executing = true
					u.state.buffer = ""
					u.components.prompt.SetValue("")
					return u, tea.Sequence(
						promptCmd,
						u.execCommand(u.state.command),
					)
				} else {
					u.state.confirming = false
					u.state.executing = false
					u.state.buffer = ""
					u.components.prompt, promptCmd = u.components.prompt.Update(msg)
					u.components.prompt.SetValue("")
					u.components.prompt.Focus()
					if u.state.runMode == ReplMode {
						cmds = append(
							cmds,
							promptCmd,
							tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
							textinput.Blink,
						)
					} else {
						return u, tea.Sequence(
							promptCmd,
							tea.Println(fmt.Sprintf("\n%s\n", u.components.renderer.RenderWarning("[cancel]"))),
							tea.Quit,
						)
					}
				}
				u.state.command = ""
			} else {
				u.components.prompt.Focus()
				u.components.prompt, promptCmd = u.components.prompt.Update(msg)
				cmds = append(
					cmds,
					promptCmd,
					textinput.Blink,
				)
			}
		}
	// Handle AI engine execution output
	case ai.EngineExecOutput:
		var output string
		if msg.IsExecutable() {
			u.state.confirming = true
			u.state.command = msg.GetCommand()
			output = u.components.renderer.RenderContent(fmt.Sprintf("`%s`", u.state.command))
			output += fmt.Sprintf("  %s\n\n  confirm execution? [y/N]", u.components.renderer.RenderHelp(msg.GetExplanation()))
			u.components.prompt.Blur()
		} else {
			output = u.components.renderer.RenderContent(msg.GetExplanation())
			u.components.prompt.Focus()
			if u.state.runMode == CliMode {
				return u, tea.Sequence(
					tea.Println(output),
					tea.Quit,
				)
			}
		}
		u.components.prompt, promptCmd = u.components.prompt.Update(msg)
		return u, tea.Sequence(
			promptCmd,
			textinput.Blink,
			tea.Println(output),
		)
	// Handle AI engine chat stream output
	case ai.EngineChatStreamOutput:
		if msg.IsLast() {
			output := u.components.renderer.RenderContent(u.state.buffer)
			u.state.buffer = ""
			u.components.prompt.Focus()
			if u.state.runMode == CliMode {
				return u, tea.Sequence(
					tea.Println(output),
					tea.Quit,
				)
			} else {
				return u, tea.Sequence(
					tea.Println(output),
					textinput.Blink,
				)
			}
		} else {
			return u, u.awaitChatStream()
		}
	// Handle runner feedback (success, fail feedback msg)
	case run.RunOutput:
		u.state.querying = false
		u.components.prompt, promptCmd = u.components.prompt.Update(msg)
		u.components.prompt.Focus()
		//we get the success msg
		output := u.components.renderer.RenderSuccess(fmt.Sprintf("\n%s\n", msg.GetSuccessMessage()))
		if msg.HasError() {
			//getting the error msg if there's an error
			output = u.components.renderer.RenderError(fmt.Sprintf("\n%s\n", msg.GetErrorMessage()))
		}
		if u.state.runMode == CliMode {
			return u, tea.Sequence(
				tea.Println(output),
				tea.Quit,
			)
		} else {
			return u, tea.Sequence(
				tea.Println(output),
				promptCmd,
				textinput.Blink,
			)
		}
	// Handle errors
	case error:
		u.state.error = msg
		return u, nil
	}

	return u, tea.Batch(cmds...)
}

// View returns the string representation of the user interface.
// It renders different views based on the state of the UI.
func (u *Ui) View() string {
	//renders error, content depending on the state of the UI
	if u.state.error != nil {
		// Render error message
		return u.components.renderer.RenderError(fmt.Sprintf("[error] %s", u.state.error))
	}
//if you are in configuring state (defined in the struct Uistate on top), then we enter this condition
	if u.state.configuring {
		// Render configuration view
		return fmt.Sprintf(
			"%s\n%s",
			u.components.renderer.RenderContent(u.state.buffer),
			u.components.prompt.View(),
		)
	}
//querying, confirming, executing are UI states defined in the struct on top of this file
	if !u.state.querying && !u.state.confirming && !u.state.executing {
		// Render prompt view
		return u.components.prompt.View()
	}

	if u.state.promptMode == ChatPromptMode {
		// if we are in the promptMode state, we will render chat mode view
		// Render chat mode view
		return u.components.renderer.RenderContent(u.state.buffer)
	} else {
		if u.state.querying {
			// Call spinner.View() function from spinner.go file
			return u.components.spinner.View()
		} else {
			if !u.state.executing {
				// Render content view
				return u.components.renderer.RenderContent(u.state.buffer)
			}
		}
	}

	return ""
}

// startRepl is a method of the Ui struct that starts the REPL (Read-Eval-Print Loop) mode.
func (u *Ui) startRepl(config *config.Config) tea.Cmd {
	return tea.Sequence(
		tea.ClearScreen,
		tea.Println(u.components.renderer.RenderContent(u.components.renderer.RenderHelpMessage())),
		textinput.Blink,
		func() tea.Msg {
			u.config = config

			// Set the prompt mode based on the default prompt mode in the configuration
			if u.state.promptMode == DefaultPromptMode {
				u.state.promptMode = GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
			}

			engineMode := ai.ExecEngineMode
			if u.state.promptMode == ChatPromptMode {
				engineMode = ai.ChatEngineMode
			}

			// Create a new engine with the specified engine mode and configuration
			engine, err := ai.NewEngine(engineMode, config)
			if err != nil {
				return err
			}

			if u.state.pipe != "" {
				engine.SetPipe(u.state.pipe)
			}

			u.engine = engine
			u.state.buffer = "Welcome \n\n"
			u.state.command = ""
			u.components.prompt = NewPrompt(u.state.promptMode)

			return nil
		},
	)
}

// startCli is a method of the Ui struct that starts the CLI (Command Line Interface) mode.
// It initializes the engine, sets the prompt mode, and handles different modes of execution.
func (u *Ui) startCli(config *config.Config) tea.Cmd {
	u.config = config

	// Set the prompt mode based on the default prompt mode in the configuration
	if u.state.promptMode == DefaultPromptMode {
		u.state.promptMode = GetPromptModeFromString(config.GetUserConfig().GetDefaultPromptMode())
	}

	engineMode := ai.ExecEngineMode
	if u.state.promptMode == ChatPromptMode {
		engineMode = ai.ChatEngineMode
	}

	// Create a new engine with the specified engine mode and configuration
	engine, err := ai.NewEngine(engineMode, config)
	if err != nil {
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine
	u.state.querying = true
	u.state.confirming = false
	u.state.buffer = ""
	u.state.command = ""

	if u.state.promptMode == ExecPromptMode {
		// If the prompt mode is ExecPromptMode, execute the completion command
		return tea.Batch(
			u.components.spinner.Tick,
			func() tea.Msg {
				output, err := u.engine.ExecCompletion(u.state.args)
				u.state.querying = false
				if err != nil {
					return err
				}

				return *output
			},
		)
	} else {
		// If the prompt mode is ChatPromptMode, start the chat stream and await the response
		return tea.Batch(
			u.startChatStream(u.state.args),
			u.awaitChatStream(),
		)
	}
}

// startConfig is a method of the Ui struct that starts the configuration mode.
//when user has not provided the config file and we're placing the terminal into
//REPL mode, we then call this function to accept config values from user while chatting
func (u *Ui) startConfig() tea.Cmd {
	return func() tea.Msg {
		// Set the UI state
		u.state.configuring = true
		u.state.querying = false
		u.state.confirming = false
		u.state.executing = false

		// Update the buffer with the rendered configuration message
		u.state.buffer = u.components.renderer.RenderConfigMessage()
		u.state.command = ""

		// Initialize a new prompt with ConfigPromptMode
		//call the newprompt function with configPromptMode
		u.components.prompt = NewPrompt(ConfigPromptMode) // this function is in prompt.go

		return nil
	}
}

// finishConfig is a method of the Ui struct that finishes the configuration process.
func (u *Ui) finishConfig(key string) tea.Cmd {
	// Update UI state
	u.state.configuring = false

	// Write configuration to file
	config, err := config.WriteConfig(key, true)
	if err != nil {
		u.state.error = err
		return nil
	}

	u.config = config

	// Initialize AI engine
	engine, err := ai.NewEngine(ai.ExecEngineMode, config)
	if err != nil {
		u.state.error = err
		return nil
	}

	if u.state.pipe != "" {
		engine.SetPipe(u.state.pipe)
	}

	u.engine = engine

	if u.state.runMode == ReplMode {
		// If in REPL mode, return a sequence of commands
		return tea.Sequence(
			tea.ClearScreen,
			tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]\n")),
			textinput.Blink,
			func() tea.Msg {
				u.state.buffer = ""
				u.state.command = ""
				u.components.prompt = NewPrompt(ExecPromptMode)

				return nil
			},
		)
	} else {
		if u.state.promptMode == ExecPromptMode {
			// If in CLI mode with ExecPromptMode, return a sequence of commands
			u.state.querying = true
			u.state.configuring = false
			u.state.buffer = ""
			return tea.Sequence(
				tea.Println(u.components.renderer.RenderSuccess("\n[settings ok]")),
				u.components.spinner.Tick,
				func() tea.Msg {
					output, err := u.engine.ExecCompletion(u.state.args)
					u.state.querying = false
					if err != nil {
						return err
					}

					return *output
				},
			)
		} else {
			// If in CLI mode with ChatPromptMode, return a batch of commands
			return tea.Batch(
				u.startChatStream(u.state.args),
				u.awaitChatStream(),
			)
		}
	}
}

// startExec is a method of the Ui struct that starts the execution of a command.
func (u *Ui) startExec(input string) tea.Cmd {
	return func() tea.Msg {
		u.state.querying = true
		u.state.confirming = false
		u.state.buffer = ""
		u.state.command = ""

		output, err := u.engine.ExecCompletion(input)
		u.state.querying = false
		if err != nil {
			return err
		}

		return *output
	}
}

// startChatStream is a method of the Ui struct that starts the chat stream.
func (u *Ui) startChatStream(input string) tea.Cmd {
	return func() tea.Msg {
		u.state.querying = true
		u.state.executing = false
		u.state.confirming = false
		u.state.buffer = ""
		u.state.command = ""

		err := u.engine.ChatStreamCompletion(input)
		if err != nil {
			return err
		}

		return nil
	}
}

// awaitChatStream is a method of the Ui struct that awaits the chat stream response.
func (u *Ui) awaitChatStream() tea.Cmd {
	return func() tea.Msg {
		output := <-u.engine.GetChannel()
		u.state.buffer += output.GetContent()
		u.state.querying = !output.IsLast()

		return output
	}
}

// execCommand is a method of the Ui struct that executes a command.
func (u *Ui) execCommand(input string) tea.Cmd {
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = true

	c := run.PrepareInteractiveCommand(input)

	return tea.ExecProcess(c, func(error error) tea.Msg {
		u.state.executing = false
		u.state.command = ""

		return run.NewRunOutput(error, "[error]", "[ok]")
	})
}

// editSettings is a method of the Ui struct that handles editing the settings.
func (u *Ui) editSettings() tea.Cmd {
	// Update UI state
	u.state.querying = false
	u.state.confirming = false
	u.state.executing = true

	// Prepare and execute the edit settings command
	c := run.PrepareEditSettingsCommand(fmt.Sprintf(
		"%s %s",
		u.config.GetSystemConfig().GetEditor(),
		u.config.GetSystemConfig().GetConfigFile(),
	))

	return tea.ExecProcess(c, func(error error) tea.Msg {
		// Update UI state
		u.state.executing = false
		u.state.command = ""

		if error != nil {
			// Handle error output
			return run.NewRunOutput(error, "[settings error]", "")
		}

		// Create a new config instance
		config, error := config.NewConfig()
		if error != nil {
			// Handle error output
			return run.NewRunOutput(error, "[settings error]", "")
		}

		// Update UI config and engine
		u.config = config
		engine, error := ai.NewEngine(ai.ExecEngineMode, config)
		if u.state.pipe != "" {
			engine.SetPipe(u.state.pipe)
		}
		if error != nil {
			// Handle error output
			return run.NewRunOutput(error, "[settings error]", "")
		}
		u.engine = engine

		// Return success output
		return run.NewRunOutput(nil, "", "[settings ok]")
	})
}
