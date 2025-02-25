package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Results struct {
	// height     int
	width       int
	text        string
	textColor   string
	borderColor string
	bgColor     string
	options     []string
	horizontal  string
	vertical    string
	selected    map[string]bool
}

func (r Results) Init() tea.Cmd {
	return nil
}

func (r Results) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return r, nil
}

func (r Results) View() string {
	return r.getStyle().Render(r.text)
}

func (r *Results) getStyle() lipgloss.Style {
	if len(r.selected) > 0 {
		alignHorizontal := getAlignment(r.horizontal)
		alignVertical := getAlignment(r.vertical)
		return lipgloss.NewStyle().
			Width((r.width/3)-6).
			Height(20).
			Padding(1, 2).
			Background(lipgloss.Color(r.bgColor)).
			Foreground(lipgloss.Color(r.textColor)).
			BorderForeground(lipgloss.Color(r.borderColor)).
			BorderBackground(lipgloss.Color(r.bgColor)).
			Bold(r.selected[r.options[0]]).
			Italic(r.selected[r.options[1]]).
			Strikethrough(r.selected[r.options[2]]).
			Underline(r.selected[r.options[3]]).
			Faint(r.selected[r.options[4]]).
			Reverse(r.selected[r.options[5]]).
			AlignHorizontal(alignHorizontal).
			AlignVertical(alignVertical).
			Border(lipgloss.RoundedBorder(), r.selected[r.options[6]])
	}
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.HiddenBorder()).
		Height(20).
		Width((r.width / 3) - 6)
}

func getAlignment(alignment string) lipgloss.Position {
	var align lipgloss.Position
	var positionMap = map[string]lipgloss.Position{
		"Left":   lipgloss.Left,
		"Right":  lipgloss.Right,
		"Center": lipgloss.Center,
		"Top":    lipgloss.Top,
		"Bottom": lipgloss.Bottom,
	}
	if val, ok := positionMap[alignment]; ok {
		align = val
	}
	return align
}
