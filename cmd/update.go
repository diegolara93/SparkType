package cmd

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.ChosenView {
	case 0:
		return m.homeUpdate(msg)
	case 1:
		return m.typerUpdate(msg)
	case 2:
		return m.settingsUpdate(msg)
	case 3:
		return m.recordsUpdate(msg)
	case 4:
		return m.gameOverUpdate(msg)
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
		if m.startedTyping.IsZero() {
			m.startedTyping = time.Now()
		}
		switch msg.String() { // change this to be the typing view update

		case "ctrl+c":
			return m, tea.Quit
		}
		/*
			If the length of the typed keys equals the length of the keys and the last key is the correct key, quit,
			TODO: add a popup after finishing showing wpm, accuracy, etc.
		*/

		if len(m.TypedKeys) == len(m.Keys) && m.TypedKeys[len(m.Keys)-1] == m.Keys[len(m.Keys)-1] {
			m.ChosenView = 4
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

		// deals with line wrapping
		// Where the next line starts to not break the user interface, even if the user types a random character
		if next == '\n' {
			m.TypedKeys = append(m.TypedKeys, next)
			if char == ' ' {
				return m, nil
			}
		}
		if len(m.TypedKeys) >= 1 {
			m.wpm = (m.score / charsPerWord) / (time.Since(m.startedTyping).Minutes())
		}
		m.TypedKeys = append(m.TypedKeys, msg.Runes...)
		if char == next {
			m.score += 1.
		}
	case tea.WindowSizeMsg:
		m.homeScreenMarginLR = msg.Width
		m.homeScreenMarginUB = msg.Height / 8
		return m, cmd

	case tickMsg:
		if len(m.TypedKeys) >= 1 {
			m.timeRemaining -= 1
			if m.timeRemaining == 0 {
				m.ChosenView = 4
				m.TypedKeys = []rune{} // Clear the typed keys after finishing
				return m, nil
			}
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		} else {
			return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	}
	return m, nil
}

func (m Model) recordsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) settingsUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: add stuff here for the settings update
	// var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) gameOverUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.ChosenView = 0
			return m, nil
		}
	}
	return m, nil
}
