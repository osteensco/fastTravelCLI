package settings

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ViewFunc func(tea.Model) string

func renderSettingsView(m settingsModel) string {
	// TODO add ascii art at the top centered between columns
	m.list.Title = "Settings"

	width := min(70, m.width/2)
	// height := m.docStyle.GetHeight() // this doesnt actually do anything
	
	style := lg.NewStyle().Width(width).Padding(1)

	right := ""
	left := m.list.View()

	if m.selectedView != nil && m.selectedModel != nil {
		right = m.selectedView(m)
	}
	if m.focusDetail {
		left = style.Border(lg.HiddenBorder()).Render(left)
		right = style.Height(lg.Height(left)-2).Border(lg.RoundedBorder()).Render(right)
	} else {
		left = style.Border(lg.RoundedBorder()).Render(left)
		right = style.Height(lg.Height(left)-2).Border(lg.HiddenBorder()).Render(right)
	}
	
	return lg.JoinHorizontal(lg.Top,left,right)
}

func renderCascadeOrderView(model tea.Model) string {
	mm, ok := model.(settingsModel)
	m, ok := mm.selectedModel.(cascadeModel)
	if !ok {
		panic("Wrong model passed to renderCascadeOrderView")
	}
	style := m.delegate.Styles.NormalDesc
	title := m.delegate.Styles.NormalTitle.Render("fastTravelCLI - Cascade Path Query Order")
	help := "Some help stuff here..."

	// s.NormalTitle = lg.NewStyle().
	// 	Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
	// 	Padding(0, 0, 0, 2) //nolint:mnd
	//
	// s.NormalDesc = s.NormalTitle.
	// 	Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	selectedStyle := lg.NewStyle().
		BorderForeground(lg.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}).
		Foreground(lg.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})


	var b strings.Builder
	for i, item := range m.items {
		premark := " "
		postmark := " "
		if i == m.cursor {
			premark = ">"
			postmark = "<"
			style = selectedStyle
		}

		if i == m.selected {
			premark = "|"
			postmark = "|"
			style = selectedStyle
		}

		renderedItem := fmt.Sprintf("%s %s %s\n", premark, item, postmark)

		b.WriteString(style.Render(renderedItem))
	}

	return lg.JoinVertical(lg.Left, title, b.String(), help)
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
