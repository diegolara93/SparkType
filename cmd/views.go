package cmd

import (
	"fmt"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/common-nighthawk/go-figure"
	"github.com/muesli/gamut"
)

func (m Model) View() string { // TODO: Make this a switch statement else-ifs so ugly
	if m.ChosenView == 0 { // Home view
		s := homeView(m)
		return s
	} else if m.ChosenView == 1 { // Typing game view
		s, _ := typeView(m)
		return s
	} else if m.ChosenView == 2 { // Settings view
		return m.settingsView()
	} else if m.ChosenView == 3 {
		return m.recordsView()
	} else if m.ChosenView == 4 {
		return m.gameOverView()
	}
	return "ERROR: unknown view	"
}

func homeView(m Model) string {
	blends := gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 75)
	borderStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder()) //.Margin(m.homeScreenMarginUB, m.homeScreenMarginLR/4, m.homeScreenMarginUB/4, m.homeScreenMarginLR/3).Border(lipgloss.NormalBorder())
	asciiFigure := figure.NewFigure("SparkType!", "", true)
	asciiFigureRainbow := rainbow(textStyle, asciiFigure.String(), blends)
	// tempStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFFF")).Bold(true)
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#940101")).Bold(true)
	// info := tempStyle.Render("A terminal-based typing game made with Go")
	// g := textStyle.Render("Start Typing!")
	// t := textStyle.Render("My Records")
	warning := warningStyle.Render("WARNING: Avoid constantly resizing terminal for best experience.")
	// settings := textStyle.Render("Settings")

	n := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB, lipgloss.Center, lipgloss.Top, borderStyle.Render(asciiFigureRainbow))
	// view := lipgloss.JoinVertical(lipgloss.Center, n, info, g, t, settings, warning)
	view := lipgloss.JoinVertical(lipgloss.Center, n, m.viewList.View(), warning)
	return view
}

func (m Model) recordsView() string {
	return "Here is your top 10 highest scores!"
}

func typeView(m Model) (string, tea.Model) {
	remaining := m.Keys[len(m.TypedKeys):]
	timeRemaining := textStyle.Render(strconv.Itoa(int(m.timeRemaining)))
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
	if len(m.TypedKeys) >= 1 {
		m.wpm = (m.score / charsPerWord) / (time.Since(m.startedTyping).Minutes())
	}
	wpmText := textStyle.Render(strconv.FormatFloat(m.wpm, 'f', 0, 64))
	text := textBox.Render(ansi.Wordwrap(s, 120, "\n"))
	textBox := lipgloss.JoinVertical(lipgloss.Center, text, timeRemaining, wpmText)
	textView := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB*5, lipgloss.Center, lipgloss.Center, textBox)
	return textView, m
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

func (m Model) gameOverView() string {
	m.ChosenView = 4

	results := fmt.Sprintf("Finished!\n"+
		"WPM: %.1f\n"+
		"Press Enter To Go Home", m.wpm)
	text := textBox.Render(results)

	centeredText := lipgloss.Place(m.homeScreenMarginLR, m.homeScreenMarginUB*8, lipgloss.Center, lipgloss.Center, text)
	return centeredText
}
