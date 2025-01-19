package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
	"image/color"
)

var textStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#EDFF82"))

type Model struct {
	TextInput          textinput.Model
	Words              []rune
	ChosenView         int32
	homeScreenMarginLR int
	homeScreenMarginUB int
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "cat"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return Model{
		TextInput:          ti,
		Words:              []rune{'a', 'b', 'c'},
		ChosenView:         1,
		homeScreenMarginLR: 0,
		homeScreenMarginUB: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.TextInput.Width = msg.Width
		m.homeScreenMarginLR = msg.Width / 3
		m.homeScreenMarginUB = msg.Height / 16
	}
	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.ChosenView == 0 {
		s := ""
		for _, words := range m.Words {
			s += textStyle.Render(string(words)) + " "
		}
		s += "\n" + m.TextInput.View() + ""
		s += "\n"
		return s
	} else if m.ChosenView == 1 {
		s := HomeView(m)
		return s
	}
	return "view3"
}

func HomeView(m Model) string {
	blends := gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 75)
	s := lipgloss.NewStyle().Margin(m.homeScreenMarginUB, m.homeScreenMarginLR, m.homeScreenMarginUB/3, m.homeScreenMarginLR)
	j := figure.NewFigure("SparkType!", "", true)
	d := rainbow(textStyle, j.String(), blends)
	g := textStyle.Render("Start Typing!")
	t := textStyle.Render("My Records")
	settings := textStyle.Render("Settings")
	n := lipgloss.JoinVertical(lipgloss.Center, s.Render(d), g, t, settings)
	return n
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
