package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"image/color"
	"os"
)

var textStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#9c350c"))

type model struct {
	textInput textinput.Model
	words     []rune
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "cat"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return model{
		textInput: ti,
		words:     []rune{'a', 'b', 'c'},
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	blends := gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)
	j := figure.NewFigure("SparkyType!", "", true)
	s := ""
	for _, words := range m.words {
		s += textStyle.Render(string(words)) + " "
	}
	s += "\n" + m.textInput.View()
	s += "\n" + rainbow(textStyle, j.String(), blends)
	return s + "\n"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
