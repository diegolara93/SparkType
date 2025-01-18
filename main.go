package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type model struct {
	words []string
}

func initialModel() model {
	return model{
		words: []string{"cat", "dog", "sparky"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := ""
	for _, words := range m.words {
		s += words + "\n"
	}
	return s + "\n"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
