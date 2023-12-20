package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//this entire file is for creating a prompt for the open ai api, there is no processing here
//we are just formatting and showing the CLI related UI parts

const (
	exec_icon          = "🚀 > "
	exec_placeholder   = "Execute something..."
	config_icon        = "🔒 > "
	config_placeholder = "Enter your OpenAI key..."
	chat_icon          = "💬 > "
	chat_placeholder   = "Ask me something..."
)

// Prompt is a struct that represents a prompt in the user interface.
type Prompt struct {
	mode  PromptMode      // The mode of the prompt.
	input textinput.Model // The text input model of the prompt.
}

// NewPrompt is a function that creates a new Prompt instance.
func NewPrompt(mode PromptMode) *Prompt {
	input := textinput.New()
	// Set the placeholder, text style, and prompt of the text input model based on the prompt mode.
	input.Placeholder = getPromptPlaceholder(mode)
	input.TextStyle = getPromptStyle(mode)
	input.Prompt = getPromptIcon(mode)

	// If the prompt mode is configuration, set the echo mode of the text input model to password.
	if mode == ConfigPromptMode {
		input.EchoMode = textinput.EchoPassword
	}

	// Focus the text input model.
	input.Focus()

	return &Prompt{
		mode:  mode,
		input: input,
	}
}

// GetMode is a method on the Prompt struct that returns the prompt mode.
func (p *Prompt) GetMode() PromptMode {
	return p.mode
}

// SetMode is a method on the Prompt struct that sets the prompt mode and updates the text input model accordingly.
func (p *Prompt) SetMode(mode PromptMode) *Prompt {
	p.mode = mode

	p.input.TextStyle = getPromptStyle(mode)
	p.input.Prompt = getPromptIcon(mode)
	p.input.Placeholder = getPromptPlaceholder(mode)

	return p
}

// SetValue is a method on the Prompt struct that sets the value of the text input model.
func (p *Prompt) SetValue(value string) *Prompt {
	p.input.SetValue(value)

	return p
}

// GetValue is a method on the Prompt struct that returns the value of the text input model.
func (p *Prompt) GetValue() string {
	return p.input.Value()
}

// Blur is a method on the Prompt struct that unfocuses the text input model.
func (p *Prompt) Blur() *Prompt {
	p.input.Blur()

	return p
}

// Focus is a method on the Prompt struct that focuses the text input model.
func (p *Prompt) Focus() *Prompt {
	p.input.Focus()

	return p
}

// Update is a method on the Prompt struct that updates the text input model with a message.
func (p *Prompt) Update(msg tea.Msg) (*Prompt, tea.Cmd) {
	var updateCmd tea.Cmd
	p.input, updateCmd = p.input.Update(msg)

	return p, updateCmd
}

// View is a method on the Prompt struct that returns a string representation of the text input model.
func (p *Prompt) View() string {
	return p.input.View()
}

// AsString is a method on the Prompt struct that returns a string representation of the prompt.
func (p *Prompt) AsString() string {
	style := getPromptStyle(p.mode)

	return fmt.Sprintf("%s%s", style.Render(getPromptIcon(p.mode)), style.Render(p.input.Value()))
}

// getPromptStyle is a function that returns the style of the prompt based on the prompt mode.
func getPromptStyle(mode PromptMode) lipgloss.Style {
	switch mode {
	case ExecPromptMode:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(exec_color))
	case ConfigPromptMode:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(config_color))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color(chat_color))
	}
}

// getPromptIcon is a function that returns the icon of the prompt based on the prompt mode.
func getPromptIcon(mode PromptMode) string {
	style := getPromptStyle(mode)

	switch mode {
	case ExecPromptMode:
		return style.Render(exec_icon)
	case ConfigPromptMode:
		return style.Render(config_icon)
	default:
		return style.Render(chat_icon)
	}
}

// getPromptPlaceholder is a function that returns the placeholder text of the prompt based on the prompt mode.
func getPromptPlaceholder(mode PromptMode) string {
	switch mode {
	case ExecPromptMode:
		return exec_placeholder
	case ConfigPromptMode:
		return config_placeholder
	default:
		return chat_placeholder
	}
}
