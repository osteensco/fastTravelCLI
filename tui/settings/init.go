package settings

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	ftdata "github.com/osteensco/fastTravelCLI/data"
)

// func setSizeMsg(w, h int) func() tea.Msg {
// 	sizeCmd := func() tea.Msg {
// 		return tea.WindowSizeMsg{Width: w, Height: h}
// 	}
//
// 	return sizeCmd
// }

func initSettings(settings *ftdata.Settings, file *os.File) settingsModel {
	settingsOptions := []list.Item{
		settingItem{
			name:  "Query Order",
			desc:  "Determine the ordering in which fastTravelCLI resolves a query.",
			model: 0,
		},
		settingItem{
			// TODO: Rename to something like "info", "build info", "version info"
			name:  "Build Info",
			desc:  "View fastTravelCLI build information.",
			model: 1,
		},
		// settingItem{
		// 	name:  "Manage Bookmarks",
		// 	desc:  "View and manage bookmarks saved with fastTravelCLI.",
		// 	model: 2,
		// },
		//
	}

	docStyle := lg.NewStyle().Margin(10, 2)
	delegate := list.NewDefaultDelegate() // ItemDelegate interface for changing style of items

	settingsList := list.New(settingsOptions, delegate, 0, 0)
	sp := &settingsList.Styles
	sp.Title = sp.Title.Background(lg.Color("23"))
	settingsList.SetShowHelp(false)
	settingsList.SetShowPagination(false)
	settingsList.SetFilteringEnabled(false)
	settingsList.SetShowStatusBar(false)

	modelList := []DetailModel{
		newCascadeModel(settings),
		newVersionModel(),
		// &bookmarksModel{
		// 	name: "BOOKMARKS",
		// },
		// &versionModel{
		// 	name: "VERSION INFO",
		// },
	}

	keybindHints := "↑/↓ (j/k) Navigate • Enter Select • ←/→ (h/l) Change Focus • Esc Quit"

	return settingsModel{
		settings:     settings,
		settingsFile: file,
		docStyle:     docStyle,
		list:         settingsList,
		models:       modelList,
		style:        baseStyle,
		listStyle:    sp,
		keybindHints: keybindHints,
	}
}

var baseStyle = lg.NewStyle().Padding(1)
var parentModel settingsModel

func Run(settings *ftdata.Settings, file *os.File) error {
	parentModel = initSettings(settings, file)
	p := tea.NewProgram(parentModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
