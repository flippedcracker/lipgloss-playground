package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

type model struct {
	fgColor      textinput.Model
	bgColor      textinput.Model
	focused      int
	bold         bool
	italic       bool
	border       bool
	borderWidth  int
	borderHeight int
	borderStyle  string
	borderColor  string
	alignment    string
	preview      lipgloss.Style
	exportCode   string
	text         textinput.Model
	overlay      tea.Model
}

var (
	w      int
	h      int
	blends = gamut.Blends(lipgloss.Color("#F25D94"), lipgloss.Color("#EDFF82"), 50)
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter sample text"
	ti.Focus()

	fg := textinput.New()
	fg.Placeholder = "Enter foreground color (e.g., #FF5733)"

	bg := textinput.New()
	bg.Placeholder = "Enter background color (e.g., #3333FF)"

	return model{
		fgColor: fg,
		bgColor: bg,
		text:    ti,
		focused: 0,
		preview: lipgloss.NewStyle(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w = msg.Width
		h = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		case "enter":
			// Apply the styles when Enter is pressed
			fg := m.fgColor.Value()
			bg := m.bgColor.Value()

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
			if m.border {
				style = style.Border(lipgloss.NormalBorder())
			}

			m.preview = style
			m.exportCode = fmt.Sprintf(
				"lipgloss.NewStyle().Foreground(lipgloss.Color(\"%s\")).Background(lipgloss.Color(\"%s\")).Bold(%t).Italic(%t)",
				fg, bg, m.bold, m.italic,
			)
		case "ctrl+o":
			m.overlay.View()
		case "tab", "down":
			m.focused = (m.focused + 1) % 3
			m.updateFocus()
		case "shift+tab", "up":
			m.focused = (m.focused - 1 + 3) % 3
			m.updateFocus()
		case "ctrl+b":
			m.text.Blur()
			m.fgColor.Blur()
			m.bgColor.Blur()
			m.bold = !m.bold
		case "ctrl+i":
			m.text.Blur()
			m.fgColor.Blur()
			m.bgColor.Blur()
			m.italic = !m.italic

		}
	}

	if m.text.Focused() {
		m.text, cmd = m.text.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.fgColor.Focused() {
		m.fgColor, cmd = m.fgColor.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.bgColor.Focused() {
		m.bgColor, cmd = m.bgColor.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(rainbow(lipgloss.NewStyle(), "Lipgloss Playground\n\n", blends)))
	sb.WriteString(fmt.Sprintf("Sample Text: %s\n", m.text.View()))
	sb.WriteString(fmt.Sprintf("Foreground Color: %s\n", m.fgColor.View()))
	sb.WriteString(fmt.Sprintf("Background Color: %s\n", m.bgColor.View()))
	sb.WriteString(fmt.Sprintf("Border Color: %s\n", m.borderColor))
	sb.WriteString(fmt.Sprintf("Bold: %t (toggle with 'ctrl+b')\n", m.bold))
	sb.WriteString(fmt.Sprintf("Italic: %t (toggle with 'ctrl+i')\n\n", m.italic))
	sb.WriteString(fmt.Sprintf("Preview: \n%s\n\n", m.preview.Render(m.text.Value())))
	sb.WriteString(fmt.Sprintf("Generated Code:\n%s\n\n", m.exportCode))
	sb.WriteString("[Enter] Apply | [Esc] Quit\n")

	return sb.String()

}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting program: %v\n", err)
	}
}

func (m *model) updateFocus() {
	// Update focus based on the `focused` index
	m.text.Blur()
	m.fgColor.Blur()
	m.bgColor.Blur()

	switch m.focused {
	case 0:
		m.text.Focus()
	case 1:
		m.fgColor.Focus()
	case 2:
		m.bgColor.Focus()
	}
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}
