package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	ftdata "github.com/osteensco/fastTravelCLI/data"
)

type ViewFunc func(tea.Model) string

func renderSettingsView(m *settingsModel) string {
// 	header := `
//      __           _  _____                     _   ___   __   _____ - -  -  -   -   -
//     / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ - -  -  -   -   -
//    | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/  - -  -   -   -
//    |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  - -  -  -   -   -
//    |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/ - -  -  -   -   -
//
// `

	header := ftdata.Logo
	m.list.Title = "Settings"

	width := min(70, m.style.GetWidth()/2)
	height := min(10, m.style.GetHeight()/2)
	m.style = m.style.Width(width).Height(height)

	left := ""
	right := ""

	if m.selectedModel != nil {
		m.selectedModel.ShowFocus(m.focusDetail)
		right = m.selectedModel.View()
	}
	if m.focusDetail {
		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("")).Foreground(lg.Color("62"))
		left = m.list.View()
		left = m.style.Border(lg.HiddenBorder()).Render(left)
		right = m.style.Height(lg.Height(left) - 2).Border(lg.RoundedBorder()).Render(right)
	} else {
		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("62"))
		left = m.list.View()
		left = m.style.Border(lg.RoundedBorder()).Render(left)
		right = m.style.Height(lg.Height(left) - 2).Border(lg.HiddenBorder()).Render(right)
	}

	return lg.JoinVertical(lg.Top, header, lg.JoinHorizontal(lg.Center, left, right))
}

func renderCascadeOrderView(m *cascadeModel) string {

	if m.focus {
		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("62")).Foreground(lg.Color("230"))
	} else {
		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("")).Foreground(lg.Color("62"))
	}

	for i, item := range m.list.Items() {
		item, ok := item.(cascadeItem)
		if !ok {
			panic("item is not a cascadeItem!")
		}
		premark := " "
		postmark := " "

		if i == m.selected {
			premark = "|"
			postmark = "|"
		}

		renderedItem := fmt.Sprintf("%s %s %s", premark, item.name, postmark)
		m.list.SetItem(i, cascadeItem{title: renderedItem, name: item.name})
	}
	m.list.Title = "Cascade Order"
	return m.list.View()
}

func renderManageBookmarksView(m *bookmarksModel) string {
	// TODO implement

	// Table of bookmarks, selections for editing, adding, deleting
	return "bookmarks"
}

func renderVersionView(m *versionModel) string {
	// TODO implement

	// Print version and commit hash
	return "version info"
}
