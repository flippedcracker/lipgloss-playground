package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	// Define styles
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF5F87"))
	selectedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	unselectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB"))
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFBF00")).Bold(true)
	checkboxChecked   = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render("✔")
	checkboxUnchecked = lipgloss.NewStyle().Foreground(lipgloss.Color("#BBBBBB")).Render("☐")
	footerStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Italic(true)
	leftColumnStyle   = lipgloss.NewStyle()
	middleColumnStyle = lipgloss.NewStyle()
	resultsStyle      = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder())
)

type model struct {
	textInput        textinput.Model
	textColorInput   textinput.Model
	borderColorInput textinput.Model
	bgColorInput     textinput.Model
	width            int
	height           int
	horizontal       string
	vertical         string
	options          []string
	selected         map[string]bool
	cursor           int
	focus            int
	message          string
}

var options = []string{
	"Bold",
	"Italic",
	"Strikethrough",
	"Underline",
	"Faint",
	"Reverse",
	"Border",
}

var hAlignmentOptions = []string{
	"Left",
	"Center",
	"Right",
}

var vAlignmentOptions = []string{
	"Top",
	"Center",
	"Bottom",
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Lipgloss Playground"
	ti.Width = 30
	ti.Prompt = " "
	ti.Focus()

	tc := textinput.New()
	tc.Placeholder = "Text Color"
	tc.Width = 30
	tc.Prompt = " "

	bc := textinput.New()
	bc.Placeholder = "Border Color"
	bc.Width = 30
	bc.Prompt = " "

	bg := textinput.New()
	bg.Placeholder = "Background Color"
	bg.Width = 30
	bg.Prompt = " "

	return model{
		textInput:        ti,
		textColorInput:   tc,
		borderColorInput: bc,
		bgColorInput:     bg,
		options:          options,
		selected:         make(map[string]bool),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Tab):
			m.focus = (m.focus + 1) % 7
			m.cursor = 0
		case key.Matches(msg, keys.ShiftTab):
			m.focus = (m.focus - 1 + 7) % 7
			m.cursor = 0

		case key.Matches(msg, keys.Up):
			var opts []string
			switch m.focus {
			case 1:
				opts = options
			case 2:
				opts = hAlignmentOptions
			case 3:
				opts = vAlignmentOptions
			}
			if len(opts) > 0 {
				m.cursor = (m.cursor - 1 + len(opts)) % len(opts)
			}

		case key.Matches(msg, keys.Down):
			var opts []string
			switch m.focus {
			case 1:
				opts = options
			case 2:
				opts = hAlignmentOptions
			case 3:
				opts = vAlignmentOptions
			}
			if len(opts) > 0 {
				m.cursor = (m.cursor + 1) % len(opts)
			}

		case key.Matches(msg, keys.Space):
			switch m.focus {
			case 1:
				m.selected[options[m.cursor]] = !m.selected[options[m.cursor]]
			case 2:
				m.horizontal = hAlignmentOptions[m.cursor]
				for _, opt := range hAlignmentOptions {
					m.selected[opt] = (opt == m.horizontal)
				}
			case 3:
				m.vertical = vAlignmentOptions[m.cursor]
				for _, opt := range vAlignmentOptions {
					m.selected[opt] = (opt == m.vertical)
				}
				m.message = fmt.Sprintf("message: %v, %v", m.horizontal, m.vertical)
			}

		case key.Matches(msg, keys.Enter):
			selectedItems := []string{}
			for _, opt := range options {
				if m.selected[opt] {
					selectedItems = append(selectedItems, opt)
				}
			}
			lipgloss.JoinHorizontal(lipgloss.Top,
				fmt.Sprint("\n"+titleStyle.Render("Form Submitted!")),
				fmt.Sprint("Name:", selectedStyle.Render(m.textInput.Value())),
				fmt.Sprint("Selected Languages:", selectedStyle.Render(strings.Join(selectedItems, ", "))),
			)
			return m, nil

		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		}
		switch m.focus {
		case 0:
			m.textInput.Focus()
			m.textInput, cmd = m.textInput.Update(msg)
		case 4:
			m.textColorInput.Focus()
			m.textColorInput, cmd = m.textColorInput.Update(msg)
		case 5:
			m.borderColorInput.Focus()
			m.borderColorInput, cmd = m.borderColorInput.Update(msg)
		case 6:
			m.bgColorInput.Focus()
			m.bgColorInput, cmd = m.bgColorInput.Update(msg)
		}
		return m, cmd

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			zoneActions := map[string]func(){
				"textInput": func() {
					m.focus = 0
					m.textInput.Focus()
					m.textColorInput.Blur()
					m.borderColorInput.Blur()
					m.bgColorInput.Blur()
				},
				"colorText": func() {
					m.focus = 4
					m.textInput.Blur()
					m.textColorInput.Focus()
					m.borderColorInput.Blur()
					m.bgColorInput.Blur()
				},
				"colorBorder": func() {
					m.focus = 5
					m.textInput.Blur()
					m.borderColorInput.Focus()
					m.textColorInput.Blur()
					m.bgColorInput.Blur()
				},
				"colorBg": func() {
					m.focus = 6
					m.textInput.Blur()
					m.bgColorInput.Focus()
					m.textColorInput.Blur()
					m.borderColorInput.Blur()
				},
			}

			// Generic function to check any set of zones
			checkZones := func(items []string, prefix string, action func(string)) bool {
				for _, item := range items {
					zoneID := prefix + item
					if zone.Get(zoneID).InBounds(msg) {
						action(item)
						m.message = fmt.Sprintf("message: '%v'", item)
						return true
					}
				}
				return false
			}

			// Check the single-action zones
			for zoneID, action := range zoneActions {
				if zone.Get(zoneID).InBounds(msg) {
					action()
					return m, nil
				}
			}

			// Check grouped zones
			if checkZones(options, "", func(item string) {
				m.selected[item] = !m.selected[item]
			}) {
				return m, nil
			}

			if checkZones(hAlignmentOptions, "h", func(item string) {
				m.horizontal = item
			}) {
				return m, nil
			}

			if checkZones(vAlignmentOptions, "v", func(item string) {
				m.vertical = item
			}) {
				return m, nil
			}
		}
	}
	return m, nil
}

func (m model) View() string {

	header := lipgloss.Place(m.width, 2, lipgloss.Center, lipgloss.Top, "Lipgloss Playground")

	textInput := titleStyle.Render("Sample Text") + "\n\n"
	if m.focus == 0 {
		textInput += cursorStyle.Render("> ") + unselectedStyle.Render(m.textInput.View()) + "\n\n"
	} else {
		textInput += "  " + unselectedStyle.Render(zone.Mark("textInput", m.textInput.View())) + "\n\n"
	}

	// Multi-Select Options
	optionSelector := titleStyle.Render("Select Styles:") + "\n"
	for i, opt := range options {
		cursor := "  " // Default cursor
		if m.focus == 1 && m.cursor == i {
			cursor = cursorStyle.Render("> ") // Highlighted
		} else {
			cursor = "  " // No cursor
		}

		checkmark := checkboxUnchecked // Default empty box
		if m.selected[opt] {
			checkmark = checkboxChecked // Checked box
		}

		optionSelector += fmt.Sprintf("%s %s %s\n", cursor, checkmark, zone.Mark(opt, opt))
	}

	hAlignmentSelector := (titleStyle.Render("Horizontal Alignment:") + "\n")
	for i, opt := range hAlignmentOptions {
		cursor := "  "
		if m.focus == 2 && m.cursor == i {
			cursor = cursorStyle.Render("> ")
		} else {
			cursor = "  "
		}

		checkmark := checkboxUnchecked
		if m.horizontal == opt {
			checkmark = checkboxChecked
		}

		hAlignmentSelector += fmt.Sprintf("%s %s %s\n", cursor, checkmark, zone.Mark("h"+opt, opt))
	}

	vAlignmentSelector := (titleStyle.Render("Vertical Alignment:") + "\n")
	for i, opt := range vAlignmentOptions {
		cursor := "  "
		if m.focus == 3 && m.cursor == i {
			cursor = cursorStyle.Render("> ")
		} else {
			cursor = "  "
		}

		checkmark := checkboxUnchecked
		if m.vertical == opt {
			checkmark = checkboxChecked
		}

		vAlignmentSelector += fmt.Sprintf("%s %s %s\n", cursor, checkmark, zone.Mark("v"+opt, opt))
	}

	textColor := titleStyle.Render("Text Color") + "\n\n"
	if m.focus == 4 {
		textColor += cursorStyle.Render("> ") + unselectedStyle.Render(m.textColorInput.View()) + "\n\n"
	} else {
		textColor += "  " + unselectedStyle.Render(zone.Mark("colorText", m.textColorInput.View())) + "\n\n"
	}

	borderColor := titleStyle.Render("Border Color") + "\n\n"
	if m.focus == 5 {
		borderColor += cursorStyle.Render("> ") + unselectedStyle.Render(m.borderColorInput.View()) + "\n\n"
	} else {
		borderColor += "  " + unselectedStyle.Render(zone.Mark("colorBorder", m.borderColorInput.View())) + "\n\n"
	}

	bgColor := titleStyle.Render("Background Color") + "\n\n"
	if m.focus == 6 {
		bgColor += cursorStyle.Render("> ") + unselectedStyle.Render(m.bgColorInput.View()) + "\n\n"
	} else {
		bgColor += "  " + unselectedStyle.Render(zone.Mark("colorBg", m.bgColorInput.View())) + "\n\n"
	}

	footer := lipgloss.Place(m.width, 2, lipgloss.Center, lipgloss.Bottom, footerStyle.Render("[Tab] Switch | [Space] Select | [Enter] Submit | [Esc] Quit"))

	leftColumn := zone.Mark("left", leftColumnStyle.Width(m.width/3).Render(lipgloss.JoinVertical(lipgloss.Left, textInput, optionSelector, hAlignmentSelector, vAlignmentSelector)))
	leftHeight := strings.Count(leftColumn, "\n") + 1
	middleColumn := middleColumnStyle.Width(m.width / 3).Height(m.height / 2).Render(lipgloss.JoinVertical(lipgloss.Center, "Middle Column", m.message, textColor, borderColor, bgColor))
	rightColumn := resultsStyle.Width(m.width / 3).Height(leftHeight).Render(Results{
		width:       m.width,
		text:        m.textInput.Value(),
		textColor:   m.textColorInput.Value(),
		borderColor: m.borderColorInput.Value(),
		bgColor:     m.bgColorInput.Value(),
		options:     options,
		selected:    m.selected,
		horizontal:  m.horizontal,
		vertical:    m.vertical,
	}.View())

	mainBody := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, middleColumn, rightColumn)

	return zone.Scan(lipgloss.NewStyle().Height(m.height).Render((lipgloss.JoinVertical(lipgloss.Top, header, mainBody, footer))))
}

func main() {
	zone.NewGlobal()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
