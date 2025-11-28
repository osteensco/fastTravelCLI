package settings

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

func setSizeMsg(w, h int) func() tea.Msg {
	sizeCmd := func() tea.Msg {
		return tea.WindowSizeMsg{Width: w, Height: h}
	}

	return sizeCmd
}

func initSettings() settingsModel {
	docStyle := lg.NewStyle().Margin(10, 2)
	delegate := list.NewDefaultDelegate() // ItemDelegate interface for changing style of items

	settingsList := list.New(settingsOptions, delegate, 0, 0)
	sp := &settingsList.Styles
	sp.Title = sp.Title.Background(lg.Color("23"))
	settingsList.SetShowHelp(false)
	settingsList.SetShowPagination(false)
	settingsList.SetFilteringEnabled(false)
	settingsList.SetShowStatusBar(false)

	return settingsModel{docStyle: docStyle, list: settingsList, models: modelList, style: baseStyle, listStyle: sp}
}

var baseStyle = lg.NewStyle().Padding(1)
var parentModel = initSettings()

func Run() error {
	p := tea.NewProgram(parentModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
