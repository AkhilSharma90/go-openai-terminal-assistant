package ui

//COMPLETE

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)
//definined different colors for different types of msgs
// Colors used for rendering different types of content.
const (
	exec_color    = "#ffa657"
	config_color  = "#ffffff"
	chat_color    = "#66b3ff"
	help_color    = "#aaaaaa"
	error_color   = "#cc3333"
	warning_color = "#ffcc00"
	success_color = "#46b946"
)

// Renderer is a struct that represents a renderer for different types of content.
//success, help, warning, error will all get their own different styles
//they're all defined as Style struct that exists in the lipgloss library
type Renderer struct {
	contentRenderer *glamour.TermRenderer
	successRenderer lipgloss.Style
	warningRenderer lipgloss.Style
	errorRenderer   lipgloss.Style
	helpRenderer    lipgloss.Style
}

// NewRenderer is a function that creates a new Renderer instance.
//returns an instance of Renderer struct defined above
func NewRenderer(options ...glamour.TermRendererOption) *Renderer {
	// Create a new terminal renderer with the provided options.
	contentRenderer, err := glamour.NewTermRenderer(options...)
	if err != nil {
		return nil
	}

	// Create new styles for rendering success, warning, error, and help messages.
	successRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(success_color))
	warningRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(warning_color))
	errorRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(error_color))
	helpRenderer := lipgloss.NewStyle().Foreground(lipgloss.Color(help_color)).Italic(true)

	return &Renderer{
		contentRenderer: contentRenderer,
		successRenderer: successRenderer,
		warningRenderer: warningRenderer,
		errorRenderer:   errorRenderer,
		helpRenderer:    helpRenderer,
	}
}

//five helper functions below to render different types of msgs defined in the renderer struct
// RenderContent is a method on the Renderer struct that renders general content.
func (r *Renderer) RenderContent(in string) string {
	out, _ := r.contentRenderer.Render(in)

	return out
}

// Renders a success message.
func (r *Renderer) RenderSuccess(in string) string {
	return r.successRenderer.Render(in)
}

// Renders a warning message.
func (r *Renderer) RenderWarning(in string) string {
	return r.warningRenderer.Render(in)
}

// Renders an error message.
func (r *Renderer) RenderError(in string) string {
	return r.errorRenderer.Render(in)
}

// Renders a help message.
func (r *Renderer) RenderHelp(in string) string {
	return r.helpRenderer.Render(in)
}

// RenderConfigMessage is a method on the Renderer struct that renders a configuration message.
func (r *Renderer) RenderConfigMessage() string {
	welcome := "Welcome! 👋  \n\n"
	welcome += "I cannot find a configuration file, please enter an `OpenAI API key` "
	welcome += "from https://platform.openai.com/account/api-keys so I can generate it for you."

	return welcome
}

//these are all the ways you can interact with the terminal ai assistant, displayed in help msg
// RenderHelpMessage is a method on the Renderer struct that renders a help message.
func (r *Renderer) RenderHelpMessage() string {
	help := "**Help**\n"
	help += "- `↑`/`↓` : navigate in history\n"
	help += "- `tab`   : switch between `🚀 exec` and `💬 chat` prompt modes\n"
	help += "- `ctrl+h`: show help\n"
	help += "- `ctrl+s`: edit settings\n"
	help += "- `ctrl+r`: clear terminal and reset discussion history\n"
	help += "- `ctrl+l`: clear terminal but keep discussion history\n"
	help += "- `ctrl+c`: exit or interrupt command execution\n"

	return help
}
