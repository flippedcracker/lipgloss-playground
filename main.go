package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	foreground textinput.Model
	background textinput.Model
	bold       bool
	italic     bool
	focused    int
	preview    lipgloss.Style
	exportCode string
	text       textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter sample text"
	ti.Focus()

	fg := textinput.New()
	fg.Placeholder = "Enter foreground color (e.g., #FF5733)"

	bg := textinput.New()
	bg.Placeholder = "Enter background color (e.g., #3333FF)"

	return model{
		foreground: fg,
		background: bg,
		text:       ti,
		focused:    0,
		preview:    lipgloss.NewStyle(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "enter":
			// Apply the styles when Enter is pressed
			fg := m.foreground.Value()
			bg := m.background.Value()

			style := lipgloss.NewStyle()
			if fg != "" {
				style = style.Foreground(lipgloss.Color(fg))
			}
			if bg != "" {
				style = style.Background(lipgloss.Color(bg))
			}
			if m.bold {
				style = style.Bold(true)
			}
			if m.italic {
				style = style.Italic(true)
			}

			m.preview = style
			m.exportCode = fmt.Sprintf(
				"lipgloss.NewStyle().Foreground(lipgloss.Color(\"%s\")).Background(lipgloss.Color(\"%s\")).Bold(%t).Italic(%t)",
				fg, bg, m.bold, m.italic,
			)
		case "tab", "down":
			m.focused = (m.focused + 1) % 3
			m.updateFocus()
		case "shift+tab", "up":
			m.focused = (m.focused - 1 + 3) % 3
			m.updateFocus()
		case "ctrl+b":
			m.text.Blur()
			m.foreground.Blur()
			m.background.Blur()
			m.bold = !m.bold
		case "ctrl+i":
			m.text.Blur()
			m.foreground.Blur()
			m.background.Blur()
			m.italic = !m.italic

		}
	}

	if m.text.Focused() {
		m.text, cmd = m.text.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.foreground.Focused() {
		m.foreground, cmd = m.foreground.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.background.Focused() {
		m.background, cmd = m.background.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"Style Generator\n\n"+
			"Sample Text: %s\n"+
			"Foreground Color: %s\n"+
			"Background Color: %s\n"+
			"Bold: %t (toggle with 'ctrl+b')\n"+
			"Italic: %t (toggle with 'ctrl+i')\n\n"+
			"Preview: %s\n\n"+
			"Generated Code:\n%s\n\n"+
			"[Enter] Apply | [Esc] Quit\n",
		m.text.View(),
		m.foreground.View(),
		m.background.View(),
		m.bold,
		m.italic,
		m.preview.Render(m.text.Value()),
		m.exportCode,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
}

func (m *model) updateFocus() {
	// Update focus based on the `focused` index
	m.text.Blur()
	m.foreground.Blur()
	m.background.Blur()

	switch m.focused {
	case 0:
		m.text.Focus()
	case 1:
		m.foreground.Focus()
	case 2:
		m.background.Focus()
	}
}
