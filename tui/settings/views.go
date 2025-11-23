package settings

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewFunc func(tea.Model) string

func renderSettingsView(m settingsModel) string {
	m.list.Title = "fastTravelCLI Settings"

	return m.docStyle.Render(m.list.View())
}

func renderCascadeOrderView(mm tea.Model) string {
	m, ok := mm.(cascadeModel)
	if !ok {
		panic("Wrong model passed to renderCascadeOrderView")
	}
	var b strings.Builder
	for i, item:= range m.items {
		premark := " "
		postmark := " "
		if i == m.cursor {
			premark = ">"
			postmark = "<"
		}

		if i == m.selected {
			premark = "|"
			postmark = "|"
		}

		b.WriteString(fmt.Sprintf("%s %s %s\n", premark, item, postmark))
	}
	// TODO: add footer similar to bubbletea list-default
	return b.String()
}

func renderManageBookmarksView(mm tea.Model) string {
	m, ok := mm.(bookmarksModel)
	if !ok {
		panic("Wrong model passed to renderManageBookmarksView")
	}
	// TODO implement
	return m.name
}

func renderVersionView(mm tea.Model) string {
	m, ok := mm.(versionModel)
	if !ok {
		panic("Wrong model passed to renderVersionView")
	}
	// TODO implement
	return m.name
}
