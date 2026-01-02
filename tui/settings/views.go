package settings

import (
	"fmt"

	lg "github.com/charmbracelet/lipgloss"
	ftdata "github.com/osteensco/fastTravelCLI/data"
)

//helpers

func joinWithSpacer(width int, left string, right string) string {
	style := lg.NewStyle().Width(width - lg.Width(left) - lg.Width(right) - 2)
	return lg.JoinHorizontal(lg.Center, left, style.Render(" "), right)
}

func renderHints(style lg.Style, text string) string {
	return style.Foreground(lg.Color("240")).Italic(true).Render(text)
}

// View functions

func renderSettingsView(m *settingsModel) string {

	header := ftdata.Logo
	m.list.Title = "Settings"

	width := min(70, m.style.GetWidth()/2)
	height := m.style.GetHeight() / 2
	innerHeight := height - 3 // accounts for the apply changes button
	m.style = m.style.Width(width).Height(height)
	m.list.SetSize(width, innerHeight)

	left := ""
	right := ""
	basePeripherialStyle := lg.NewStyle().Border(lg.HiddenBorder())
	keybindHints := renderHints(basePeripherialStyle, m.keybindHints)
	applyButtonStyle := basePeripherialStyle
	applyButtonText := ""

	if m.updated {
		applyButtonStyle = applyButtonStyle.Border(lg.RoundedBorder())
		applyButtonText = "Apply Changes"
	}

	if m.showAppliedMsg {
		applyButtonText = "Settings updated!"
	}

	if m.selectedModel != nil {
		m.selectedModel.ShowFocus(m.focusDetail)
		m.selectedModel.SetSize(width, height)
		right = m.selectedModel.View()
	}
	if m.focusDetail {
		keybindHints = renderHints(basePeripherialStyle, m.selectedModel.GetKeybindHints())
		// to get the style applied to the list title list.View() has to be called after Styles are updated
		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("")).Foreground(lg.Color("62"))
		left = m.list.View()
		applyButton := joinWithSpacer(width, "", applyButtonStyle.Render(applyButtonText))
		left = lg.JoinVertical(lg.Top, left, applyButton)
		// lg.Place gives us a bounded box
		left = lg.Place(width, innerHeight, lg.Left, lg.Top, left)
		right = lg.Place(width, height, lg.Left, lg.Top, right)

		left = m.style.Border(lg.RoundedBorder()).Render(left)
		right = m.style.Border(lg.RoundedBorder()).BorderForeground(lg.Color("62")).Render(right)
	} else {
		var leftStyle lg.Style

		if m.updated && m.applySelected {
			applyButtonStyle = applyButtonStyle.BorderForeground(lg.Color("62"))
			leftStyle = m.style.Border(lg.RoundedBorder())
			keybindHints = renderHints(basePeripherialStyle, "Enter Apply changes • ←/→ (h/l) Change Focus • ↑/↓ (j/k) Navigate • Esc Quit")
		} else {
			leftStyle = m.style.Border(lg.RoundedBorder()).BorderForeground(lg.Color("62"))
		}

		applyButton := joinWithSpacer(width, "", applyButtonStyle.Render(applyButtonText))

		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("62"))
		left = m.list.View()
		left = lg.JoinVertical(lg.Top, left, applyButton)
		left = lg.Place(width, innerHeight, lg.Left, lg.Top, left)
		right = lg.Place(width, height, lg.Left, lg.Top, right)

		left = leftStyle.Render(left)
		if m.selectedModel != nil {
			right = m.style.Border(lg.RoundedBorder()).Render(right)
		} else {
			right = m.style.Border(lg.HiddenBorder()).Render(right)
		}
	}

	return lg.JoinVertical(lg.Top, header, lg.JoinHorizontal(lg.Center, left, right), keybindHints)
}

func renderCascadeOrderView(m *cascadeModel) string {

	m.list.SetSize(m.width, m.height)
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

		renderedItem := fmt.Sprintf("%v %s %s %s", i+1, premark, item.name, postmark)
		m.list.SetItem(i, cascadeItem{title: renderedItem, name: item.name})
	}
	m.list.Title = "Query Order"
	return m.list.View()
}

//	func renderManageBookmarksView(m *bookmarksModel) string {
//		// TODO implement
//
//		// Table of bookmarks, selections for editing, adding, deleting
//		return "bookmarks"
//	}

func renderVersionView(m *versionModel) string {
	// TODO implement

	// Print version and commit hash
	return fmt.Sprintf("version: %s", ftdata.Version)
}
