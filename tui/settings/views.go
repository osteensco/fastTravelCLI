package settings

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	ftdata "github.com/osteensco/fastTravelCLI/data"
)

type ViewFunc func(tea.Model) string

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
	applyButtonStyle := lg.NewStyle().Border(lg.HiddenBorder())
	applyButtonText := ""
	applyButtonHint := ""

	if m.updated {
		applyButtonStyle = applyButtonStyle.Border(lg.RoundedBorder())
		applyButtonText = "Apply Changes"
		applyButtonHint = applyButtonStyle.Foreground(lg.Color("240")).Border(lg.HiddenBorder()).Render("[s]")
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
		// to get the style applied to the list title list.View() has to be called after Styles are updated
		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("")).Foreground(lg.Color("62"))
		left = m.list.View()
		// TODO move this below everything and add other keymap hints
		footer := joinWithSpacer(width, applyButtonStyle.Render(applyButtonText), applyButtonHint)
		// applyButton := lg.JoinHorizontal(lg.Left,applyButtonStyle.Render(applyButtonText),applyButtonHint)
		left = lg.JoinVertical(lg.Top, left, footer)
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
		} else {
			leftStyle = m.style.Border(lg.RoundedBorder()).BorderForeground(lg.Color("62"))
		}

		footer := joinWithSpacer(width, applyButtonStyle.Render(applyButtonText), applyButtonHint)

		m.list.Styles.Title = m.list.Styles.Title.Background(lg.Color("62"))
		left = m.list.View()
		left = lg.JoinVertical(lg.Top, left, footer)
		left = lg.Place(width, innerHeight, lg.Left, lg.Top, left)
		right = lg.Place(width, height, lg.Left, lg.Top, right)

		left = leftStyle.Render(left)
		if m.selectedModel != nil {
			right = m.style.Border(lg.RoundedBorder()).Render(right)
		} else {
			right = m.style.Border(lg.HiddenBorder()).Render(right)
		}
	}

	return lg.JoinVertical(lg.Top, header, lg.JoinHorizontal(lg.Center, left, right))
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
	return fmt.Sprintf("version: %s",ftdata.Version) 
}
