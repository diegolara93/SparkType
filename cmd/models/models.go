package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/x/ansi"
	"image/color"
	"io"
	"strings"
	"time"

	"SparkType/cmd/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

var textStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#EDFF82"))

var textBox = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#9c14db"))

var typedText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#EDFF82"))

var wrongText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#8a0101")).
	Underline(true)

type typingSettings struct {
	amountOfWords int
	punctuation   bool
	numbers       bool
}
type Model struct {
	TextInput          textinput.Model
	Keys               []rune
	TypedKeys          []rune
	ChosenView         int32
	homeScreenMarginLR int
	homeScreenMarginUB int
	viewList           list.Model
	settings           typingSettings
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
		Keys:               []rune(utils.GenerateWord(20)),
		TypedKeys:          []rune{},
		ChosenView:         0,
		homeScreenMarginLR: 0,
		homeScreenMarginUB: 0,
		viewList:           viewsList,
	}
}

func (m Model) Init() tea.Cmd {

	return tea.Batch(tick(), textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.ChosenView {
	case 0:
		return m.homeUpdate(msg)
	case 1:
		return m.typerUpdate(msg)
	case 2:
		return m.settingsUpdate()
	}
	return m, nil
}

func (m Model) homeUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			_, ok := m.viewList.SelectedItem().(item)
			if ok {
				m.ChosenView = int32(m.viewList.Index() + 1)
			}
			return m, cmd
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

func (m Model) typerUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() { // change this to be the typing view update

		case "ctrl+c", "q":
			return m, tea.Quit
		}
		/*
			If the length of the typed keys equals the length of the keys and the last key is the correct key, quit,
			TODO: add a popup after finishing showing wpm, accuracy, etc.
		*/

		if len(m.TypedKeys) == len(m.Keys) && m.TypedKeys[len(m.Keys)-1] == m.Keys[len(m.Keys)-1] {
			m.ChosenView = 0
			m.TypedKeys = []rune{} // Clear the typed keys after finishing
			return m, nil
		}
		// Deleting characters
		if msg.Type == tea.KeyBackspace && len(m.TypedKeys) > 0 {
			m.TypedKeys = m.TypedKeys[:len(m.TypedKeys)-1]
			return m, nil
		}

		// CHECKING FOR TEA.KEYSPACE ALLOWS MACOS AND SSH TO DETECT SPACE, WITHOUT IT, IT DOESN'T PICK IT UP
		if msg.Type != tea.KeyRunes && msg.Type != tea.KeySpace {
			return m, nil
		}

		char := msg.Runes[0]
		next := rune(m.Keys[len(m.TypedKeys)])

		// To properly account for line wrapping we need to always insert a new line
		// Where the next line starts to not break the user interface, even if the user types a random character
		if next == '\n' {
			m.TypedKeys = append(m.TypedKeys, next)

			// Since we need to perform a line break
			// if the user types a space we should simply ignore it.
			if char == ' ' {
				return m, nil
			}
		}

		m.TypedKeys = append(m.TypedKeys, msg.Runes...)
	case tea.WindowSizeMsg:
		m.homeScreenMarginLR = msg.Width
		m.homeScreenMarginUB = msg.Height / 8
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.ChosenView == 0 { // Home view
		s := homeView(m)
		return s
	} else if m.ChosenView == 1 { // Typing game view
		s := typeView(m)
		return s
	} else if m.ChosenView == 2 { // Settings view
		return m.settingsView()
	}
	return "temp view"
}

func typeView(m Model) string {
	remaining := m.Keys[len(m.TypedKeys):]
	var typed string
	for i, c := range m.TypedKeys {
		if c == rune(m.Keys[i]) {
			typed += typedText.Render(string(c))
		} else {
			typed += wrongText.Render(string(m.Keys[i]))
		}
	}

	s := fmt.Sprintf(
		"%s",
		typed,
	)
	if len(remaining) > 0 {
		s += string(remaining[:1])
		s += string(remaining[1:])
	}
	if len(remaining) == 0 && m.TypedKeys[len(m.Keys)-1] == m.Keys[len(m.Keys)-1] {
		return homeView(m)
	}
	text := textBox.Render(ansi.Wordwrap(s, 120, "\n"))
	textBox := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB, lipgloss.Center, lipgloss.Top, text)
	return textBox
}

func homeView(m Model) string {
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

func (m Model) recordsView() string {
	return "Your highest scores:"
}

func (m Model) recordsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) settingsView() string {
	// TODO: make the selection a list have a seperate settingsUpdate()
	header := textStyle.Render("Customize the settings of the typing game here:")
	amountOfWords := textStyle.Render("Amount of Words (Minimum 15, Maximum 150): ")
	punctuation := textStyle.Render("Punctuation: ")
	numbers := textStyle.Render("Numbers: ")
	centered := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB, lipgloss.Center, lipgloss.Top+0.1, header)
	return lipgloss.JoinVertical(lipgloss.Center, centered, amountOfWords, punctuation, numbers)
}

func (m Model) settingsUpdate() (tea.Model, tea.Cmd) {
	// TODO: add stuff here for the settings update
	return m, nil
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
