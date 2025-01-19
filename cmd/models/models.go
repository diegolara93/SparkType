package models

import (
	"image/color"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
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
		m.homeScreenMarginLR = msg.Width
		m.homeScreenMarginUB = msg.Height / 5
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
	} else if m.ChosenView == 2 {
		return m.settingsView()
	}
	return "view3"
}

func HomeView(m Model) string {
	blends := gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 75)
	borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()) //.Margin(m.homeScreenMarginUB, m.homeScreenMarginLR/4, m.homeScreenMarginUB/4, m.homeScreenMarginLR/3).Border(lipgloss.NormalBorder())
	asciiFigure := figure.NewFigure("SparkType!", "", true)
	asciiFigureRainbow := rainbow(textStyle, asciiFigure.String(), blends)
	tempStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFFF")).Bold(true)
	info := tempStyle.Render("A terminal-based typing game made with Go")
	g := textStyle.Render("Start Typing!")
	t := textStyle.Render("My Records")
	settings := textStyle.Render("Settings")
	n := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB, lipgloss.Center, lipgloss.Center, borderStyle.Render(asciiFigureRainbow))
	view := lipgloss.JoinVertical(lipgloss.Center, n, info, g, t, settings)
	return view
}

func (m Model) settingsView() string {
	s := textStyle.Render("Customize the settings of the typing game here:")
	c := textStyle.Render("Periods and SemiColons: ")
	return lipgloss.JoinVertical(lipgloss.Center, s, c)
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
