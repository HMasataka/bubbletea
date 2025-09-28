package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A sample Bubble Tea app that splits the screen into 4 quadrants
// and renders a table with dummy data in each quadrant.

var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6B7280")).
			Padding(0, 1).
			Margin(0, 0)

	borderFocused = borderStyle.Copy().
			BorderForeground(lipgloss.Color("#7D56F4")).
			Bold(true)

	titleStyle = lipgloss.NewStyle().Bold(true)
)

type model struct {
	width  int
	height int

	topLeft     table.Model
	topRight    table.Model
	bottomLeft  table.Model
	bottomRight table.Model

	focused int // 0: tl, 1: tr, 2: bl, 3: br
}

func newTable(title string, cols []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(false),
	)

	// Basic table styling
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)
	// Keep default Selected style so the cursor/selection is visible when focused.
	t.SetStyles(s)

	// Store the title in the table's internal title field by
	// rendering it above the table when composing panes.
	// We'll prepend it in the pane renderer, not inside table.Model.
	_ = title // kept for clarity; used by pane headers
	return t
}

func initialModel() model {
	// Top-left: Users
	tlCols := []table.Column{{Title: "ID", Width: 4}, {Title: "Name", Width: 14}, {Title: "Role", Width: 10}}
	tlRows := []table.Row{
		{"1", "Alice", "Admin"},
		{"2", "Bob", "Editor"},
		{"3", "Charlie", "Viewer"},
		{"4", "Diana", "Admin"},
		{"5", "Evan", "Editor"},
		{"6", "Fiona", "Viewer"},
		{"7", "Gabe", "Editor"},
		{"8", "Hana", "Viewer"},
	}
	tl := newTable("Users", tlCols, tlRows)

	// Top-right: Orders
	trCols := []table.Column{{Title: "Order", Width: 7}, {Title: "Customer", Width: 14}, {Title: "Total", Width: 8}}
	trRows := []table.Row{
		{"#1001", "Alice", "$120.00"},
		{"#1002", "Bob", "$56.80"},
		{"#1003", "Charlie", "$240.50"},
		{"#1004", "Diana", "$18.99"},
		{"#1005", "Evan", "$75.20"},
		{"#1006", "Fiona", "$310.00"},
		{"#1007", "Gabe", "$49.90"},
		{"#1008", "Hana", "$88.88"},
	}
	tr := newTable("Orders", trCols, trRows)

	// Bottom-left: Products
	blCols := []table.Column{{Title: "SKU", Width: 8}, {Title: "Product", Width: 16}, {Title: "Stock", Width: 6}}
	blRows := []table.Row{
		{"A-001", "Keyboard", "42"},
		{"A-002", "Mouse", "133"},
		{"A-003", "Monitor", "12"},
		{"A-004", "Chair", "8"},
		{"A-005", "Desk", "5"},
		{"A-006", "USB Hub", "64"},
		{"A-007", "Webcam", "27"},
		{"A-008", "Headset", "35"},
	}
	bl := newTable("Products", blCols, blRows)

	// Bottom-right: Logs
	brCols := []table.Column{{Title: "Time", Width: 8}, {Title: "Level", Width: 8}, {Title: "Message", Width: 24}}
	brRows := []table.Row{
		{"12:01", "INFO", "Server started"},
		{"12:02", "WARN", "High latency"},
		{"12:03", "INFO", "User login: alice"},
		{"12:04", "ERROR", "DB timeout"},
		{"12:05", "INFO", "Retry succeeded"},
		{"12:06", "INFO", "User logout: alice"},
		{"12:07", "INFO", "Metrics flushed"},
		{"12:08", "INFO", "Healthcheck OK"},
	}
	br := newTable("Logs", brCols, brRows)

	return model{
		topLeft:     tl,
		topRight:    tr,
		bottomLeft:  bl,
		bottomRight: br,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC || msg.String() == "esc" {
			return m, tea.Quit
		}

		// Focus switching
		switch msg.Type {
		case tea.KeyTab:
			m.setFocus((m.focused + 1) % 4)
		case tea.KeyShiftTab:
			m.setFocus((m.focused + 3) % 4)
		}

		// Route the key to the focused table so it can scroll
		var cmd tea.Cmd
		switch m.focused {
		case 0:
			t, c := m.topLeft.Update(msg)
			m.topLeft = t
			cmd = c
		case 1:
			t, c := m.topRight.Update(msg)
			m.topRight = t
			cmd = c
		case 2:
			t, c := m.bottomLeft.Update(msg)
			m.bottomLeft = t
			cmd = c
		case 3:
			t, c := m.bottomRight.Update(msg)
			m.bottomRight = t
			cmd = c
		}
		return m, cmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recomputeSizes()
		return m, nil
	default:
		return m, nil
	}
}

func (m *model) recomputeSizes() {
	// Compute pane sizes. We'll split evenly with a 1-column/row gap.
	gap := 1
	// available width/height for two panes plus gap
	w := m.width
	h := m.height
	if w < 10 || h < 10 {
		return
	}

	paneW := (w - gap) / 2
	paneH := (h - gap) / 2

	// Update column widths if panes are tight. Not strictly required,
	// but helps prevent overflow in very small terminals.
	_ = paneW
	// Table height is number of rows to render (header + rows). We'll set
	// a conservative height based on paneH minus borders and padding and title.
	// border adds 2, plus an extra title line, so subtract 3 in total.
	tableHeight := paneH - 3
	if tableHeight < 3 {
		tableHeight = 3
	}
	m.topLeft.SetHeight(tableHeight)
	m.topRight.SetHeight(tableHeight)
	m.bottomLeft.SetHeight(tableHeight)
	m.bottomRight.SetHeight(tableHeight)

	// Approximate inner content width (pane width - borders - padding)
	contentWidth := paneW - 2 - 2
	if contentWidth < 10 {
		contentWidth = 10
	}
	m.topLeft.SetWidth(contentWidth)
	m.topRight.SetWidth(contentWidth)
	m.bottomLeft.SetWidth(contentWidth)
	m.bottomRight.SetWidth(contentWidth)
}

func pane(title string, content string, width, height int, focused bool) string {
	st := borderStyle
	if focused {
		st = borderFocused
	}
	boxed := st.Width(width).Height(height)
	header := titleStyle.Render(title)
	// Render title on first line, then content below; lipgloss will clamp/truncate.
	return boxed.Render(header + "\n" + content)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		// Wait for first WindowSizeMsg
		return ""
	}

	gap := 1
	paneW := (m.width - gap) / 2
	paneH := (m.height - gap) / 2

	// Compose each quadrant
	tl := pane("Users", m.topLeft.View(), paneW, paneH, m.focused == 0)
	tr := pane("Orders", m.topRight.View(), paneW, paneH, m.focused == 1)
	bl := pane("Products", m.bottomLeft.View(), paneW, paneH, m.focused == 2)
	br := pane("Logs", m.bottomRight.View(), paneW, paneH, m.focused == 3)

	// Top row: tl | tr
	top := lipgloss.JoinHorizontal(lipgloss.Left, tl, lipgloss.NewStyle().Width(gap).Render(" "), tr)
	// Bottom row: bl | br
	bottom := lipgloss.JoinHorizontal(lipgloss.Left, bl, lipgloss.NewStyle().Width(gap).Render(" "), br)

	// Stack top and bottom with a gap row
	vgap := lipgloss.NewStyle().Height(gap).Render("\n")
	return lipgloss.JoinVertical(lipgloss.Left, top, vgap, bottom)
}

func main() {
	m := initialModel()
	m.setFocus(0)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func (m *model) setFocus(i int) {
	if i < 0 || i > 3 {
		return
	}
	m.focused = i
	m.topLeft.Blur()
	m.topRight.Blur()
	m.bottomLeft.Blur()
	m.bottomRight.Blur()
	switch i {
	case 0:
		m.topLeft.Focus()
	case 1:
		m.topRight.Focus()
	case 2:
		m.bottomLeft.Focus()
	case 3:
		m.bottomRight.Focus()
	}
}
