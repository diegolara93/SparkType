package sparktype

import (
	"SparkType/cmd/models"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func Start() {
	p := tea.NewProgram(models.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
