package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	counter  int
	typedKey string
}

var _ tea.Model = model{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.typedKey = msg.String()
			m.counter++
		}
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("You typed: %s\nCounter is: %d\nCtrl+C to exit", m.typedKey, m.counter)
}

func main() {
	m := model{}
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
