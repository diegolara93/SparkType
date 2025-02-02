package cmd

import (
	"SparkType/cmd/utils"
	"fmt"
	"image/color"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
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
	Background(lipgloss.Color("#8a0101")).
	Underline(true)

type typingSettings struct {
	amountOfWords int
	punctuation   bool
	numbers       bool
	time          int
}

const charsPerWord = 5

type Model struct {
	TextInput          textinput.Model
	Keys               []rune
	TypedKeys          []rune
	ChosenView         int32
	homeScreenMarginLR int
	homeScreenMarginUB int
	viewList           list.Model
	settings           typingSettings
	timeRemaining      int32
	score              float64
	startedTyping      time.Time
	wpm                float64
	top10Scores        []float64
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

func initialSettings() typingSettings {
	return typingSettings{
		amountOfWords: 15,
		punctuation:   false,
		numbers:       false,
		time:          30,
	}
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

	settings := initialSettings()
	return Model{
		TextInput:          ti,
		Keys:               []rune(utils.GenerateWord(30)),
		TypedKeys:          []rune{},
		ChosenView:         0,
		homeScreenMarginLR: 0,
		homeScreenMarginUB: 0,
		viewList:           viewsList,
		settings:           settings,
		timeRemaining:      10,
		score:              0,
		startedTyping:      time.Now(),
		wpm:                0.,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tick(), textinput.Blink)
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
