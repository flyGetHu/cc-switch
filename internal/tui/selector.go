package tui

import (
	"fmt"
	"strings"

	"cc-switch/internal/config"
	"cc-switch/internal/provider"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			Foreground(lipgloss.Color("10")).
			Bold(true)

	currentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginTop(1)
)

type SelectorModel struct {
	providers map[string]provider.Provider
	keys      []string
	current   string
	cursor    int
	choice    string
}

func NewSelectorModel() (SelectorModel, error) {
	cfg, err := config.Load()
	if err != nil {
		return SelectorModel{}, err
	}

	keys := make([]string, 0, len(cfg.Providers))
	for k := range cfg.Providers {
		keys = append(keys, k)
	}

	return SelectorModel{
		providers: cfg.Providers,
		keys:      keys,
		current:   cfg.Current,
		cursor:    0,
	}, nil
}

func (m SelectorModel) Init() tea.Cmd {
	return nil
}

func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.keys)-1 {
				m.cursor++
			}

		case "enter", " ":
			if len(m.keys) > 0 {
				m.choice = m.keys[m.cursor]
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m SelectorModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("选择服务商"))
	b.WriteString("\n\n")

	for i, key := range m.keys {
		p := m.providers[key]

		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		name := p.Name
		if key == m.current {
			name = fmt.Sprintf("%s (当前)", name)
		}

		line := fmt.Sprintf("%s%s", cursor, name)
		if m.cursor == i {
			line = selectedStyle.Render(line)
		} else if key == m.current {
			line = currentStyle.Render("  " + name)
		} else {
			line = itemStyle.Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("\n↑/↓ 选择 • Enter 确认 • q 退出"))

	return b.String()
}

func (m SelectorModel) GetChoice() string {
	return m.choice
}

func RunSelector() (string, error) {
	m, err := NewSelectorModel()
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	return finalModel.(SelectorModel).GetChoice(), nil
}
