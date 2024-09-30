package main

import (
	"encoding/json"
	"fmt"
	"os"

	"weezel/playground/cmd/tvprogs/programinfo"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	progGuide       programinfo.Channels
	progInfoContent string
	programInfo     viewport.Model
	channelPicker   table.Model
	width           int
	height          int
	ready           bool
}

func newChannelPicker() table.Model {
	rows := []table.Row{}
	for i, channaleName := range programinfo.ChannelOrder {
		rows = append(rows, table.Row{fmt.Sprintf("%d", i+1), channaleName})
	}

	t := table.New(
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
		table.WithWidth(30),
		table.WithColumns([]table.Column{
			{Title: "Number", Width: 6},
			{Title: "Channel", Width: 20},
		}),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)
	t.SetStyles(s)

	return t
}

func newModel() *model {
	t := newChannelPicker()

	progViewPort := viewport.New(0, 0)

	f, err := os.ReadFile("prog.json")
	if err != nil {
		panic(err)
	}

	m := &model{
		channelPicker: t,
		programInfo:   progViewPort,
	}
	err = json.Unmarshal(f, &m.progGuide)
	if err != nil {
		panic(err)
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.channelPicker, cmd = m.channelPicker.Update(msg)
			selectedChannel := m.channelPicker.SelectedRow()[1]
			m.progInfoContent = m.progGuide.GetChannelWholeDay(selectedChannel)
		case "down", "j":
			m.channelPicker, cmd = m.channelPicker.Update(msg)
			selectedChannel := m.channelPicker.SelectedRow()[1]
			m.progInfoContent = m.progGuide.GetChannelWholeDay(selectedChannel)
		case " ", "enter":
			selectedChannel := m.channelPicker.SelectedRow()[1]
			m.progInfoContent = m.progGuide.GetChannelWholeDay(selectedChannel)
		case "pageDown":
			m.programInfo, cmd = m.programInfo.Update(msg)
		case "pageUp":
			m.programInfo, cmd = m.programInfo.Update(msg)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		chPickerHeight := lipgloss.Height(m.channelPicker.View())
		footerHeight := lipgloss.Height(m.programInfo.View())
		verticalMarginHeight := chPickerHeight + footerHeight

		if !m.ready {
			m.programInfo = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.programInfo.HighPerformanceRendering = false
			m.programInfo.SetContent(m.channelPicker.SelectedRow()[1])
			m.ready = true
			m.programInfo.YPosition = chPickerHeight + 1
		} else {
			m.programInfo.Width = msg.Width
			m.programInfo.Height = msg.Height - verticalMarginHeight
		}
	}

	return m, cmd
}

// View implements tea.Model.
func (m model) View() string {
	progInfoContent := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(50).
		Height(m.channelPicker.Height()*2).
		Padding(0, 0, 0, 0).
		Border(lipgloss.RoundedBorder())
	infoView := progInfoContent.Render(m.progInfoContent)

	return lipgloss.JoinHorizontal(lipgloss.Left, m.channelPicker.View(), infoView)
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
