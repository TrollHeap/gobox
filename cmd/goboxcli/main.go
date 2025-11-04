package main

import (
	"fmt"
	"os"

	display "gobox/internal/ui/display"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count int
}

func NewModel() model {
	return model{count: 0}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.count++
		}
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("Count: %d\n", m.count)
}

func main() {
	display.DisplayBatteryReport()
	p := tea.NewProgram(NewModel())
	if err := p.Start(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}
