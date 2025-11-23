package settings

import (
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
	var docStyle = lg.NewStyle().Margin(10, 2)
	return settingsModel{docStyle: docStyle, list: settingsList, models: modelList}
}

var parentModel = initSettings()

func Run() error {
	p := tea.NewProgram(parentModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
