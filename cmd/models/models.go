package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"image/color"
	"io"
	"strings"

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
	viewList           list.Model
}
type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := textStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return textStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "cat"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	views := []list.Item{
		item("Start Typing"),
		item("Settings"),
		item("My Records"),
	}
	viewsList := list.New(views, itemDelegate{}, 20, 10)
	viewsList.SetShowTitle(false)
	viewsList.SetShowHelp(false)
	viewsList.SetShowFilter(false)
	viewsList.SetShowStatusBar(false)
	viewsList.SetShowPagination(false)
	return Model{
		TextInput:          ti,
		Words:              []rune{'a', 'b', 'c'},
		ChosenView:         1,
		homeScreenMarginLR: 0,
		homeScreenMarginUB: 0,
		viewList:           viewsList,
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
		m.homeScreenMarginUB = msg.Height / 8
	}
	m.TextInput, cmd = m.TextInput.Update(msg)
	m.viewList, cmd = m.viewList.Update(msg)
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
	//tempStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFFF")).Bold(true)
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#940101")).Bold(true)
	//info := tempStyle.Render("A terminal-based typing game made with Go")
	//g := textStyle.Render("Start Typing!")
	//t := textStyle.Render("My Records")
	warning := warningStyle.Render("WARNING: Avoid constantly resizing terminal for best experience.")
	//settings := textStyle.Render("Settings")

	n := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB, lipgloss.Center, lipgloss.Top, borderStyle.Render(asciiFigureRainbow))
	//view := lipgloss.JoinVertical(lipgloss.Center, n, info, g, t, settings, warning)
	view := lipgloss.JoinVertical(lipgloss.Center, n, m.viewList.View(), warning)
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
