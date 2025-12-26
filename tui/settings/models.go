package settings

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	ftdata "github.com/osteensco/fastTravelCLI/data"
)

type exitDetail struct{}
type unfocusDetail struct{}
type handledCommand struct{}

type DetailModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
	ShowFocus(focus bool)
}

type settingsModel struct {
	settings 	  *ftdata.Settings
	docStyle      lg.Style
	list          list.Model
	models        []DetailModel
	selectedModel DetailModel
	focusDetail   bool
	style         lg.Style
	listStyle     *list.Styles
}

type settingItem struct {
	name, desc string
	model      int
}

func (i settingItem) Title() string       { return i.name }
func (i settingItem) Description() string { return i.desc }
func (i settingItem) FilterValue() string { return i.name }

var (
	modelList = []DetailModel{
		newCascadeModel(),

		// &bookmarksModel{
		// 	name: "BOOKMARKS",
		// },
		// &versionModel{
		// 	name: "VERSION INFO",
		// },
	}

	settingsOptions = []list.Item{

		// TODO add "general" settings item.

		settingItem{
			// TODO: rename to something better, more clear/intuitive
			//	- query order?
			name:  "Query Order",
			desc:  "Determine the ordering in which fastTravelCLI resolves a query.",
			model: 0,
		},

		// settingItem{
		// 	name:  "Manage Bookmarks",
		// 	desc:  "View and manage bookmarks saved with fastTravelCLI.",
		// 	model: 1,
		// },
		//
		// settingItem{
		// 	// TODO: Rename to something like "info", "build info", "version info"
		// 	name:  "Version",
		// 	desc:  "View fastTravelCLI version information.",
		// 	model: 2,
		// },
	}
)

func (m settingsModel) Init() tea.Cmd {
	return nil
}

func (m settingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.focusDetail {
		var sm tea.Model
		sm, cmd = m.selectedModel.Update(msg)
		m.selectedModel = sm.(DetailModel)

		if cmd != nil {
			msg = cmd()
			switch msg.(type) {
			case exitDetail:
				m.focusDetail = false
				m.selectedModel = nil
				// m.selectedView = nil
			case unfocusDetail:
				m.focusDetail = false
			}
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := m.docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.style = m.style.Width(msg.Width).Height(msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		// exit the program
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// select a setting
		case "enter":
			i := m.list.SelectedItem().(settingItem)
			// m.selectedView = i.view
			m.selectedModel = m.models[i.model]
			m.focusDetail = true
			return m, nil
		case "l", "right":
			if m.selectedModel != nil {
				m.focusDetail = true
			}
		}

	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// The view is just one big string that gets passed to the UI for rendering
func (m settingsModel) View() string {
	return renderSettingsView(&m)
}

type cascadeModel struct {
	Render   ViewFunc
	list     list.Model
	cursor   int
	selected int
	focus    bool
}

type cascadeItem struct {
	title string
	name  string
	// hook func()
}

func (i cascadeItem) Title() string       { return i.title }
func (i cascadeItem) Description() string { return "" }
func (i cascadeItem) FilterValue() string { return i.title }

func newCascadeList() list.Model {
	// TODO: import items from ft settings
	items := []list.Item{cascadeItem{name: "one"}, cascadeItem{name: "two"}, cascadeItem{name: "three"}, cascadeItem{name: "four"}}
	styles := list.NewDefaultItemStyles()
	styles.SelectedTitle = lg.NewStyle().Border(lg.NormalBorder(), false).Foreground(lg.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).Padding(0, 0, 0, 1)
	delegate := list.NewDefaultDelegate()
	delegate.Styles = styles
	delegate.ShowDescription = false
	l := list.New(items, delegate, 40, 20)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.InfiniteScrolling = true
	return l
}

func newCascadeModel() DetailModel {

	m := cascadeModel{
		list:     newCascadeList(),
		cursor:   0,
		selected: -1,
	}

	return &m
}

func (m *cascadeModel) ShowFocus(focus bool) {
	m.focus = focus
}

func (m *cascadeModel) Init() tea.Cmd {
	return nil
}

func (m *cascadeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.cursor = m.list.Cursor()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected != -1 {
				if m.cursor != 0 {
					selItem := m.list.Items()[m.selected]
					swapItem := m.list.Items()[m.selected-1]
					m.list.SetItem(m.selected, swapItem)
					m.list.SetItem(m.selected-1, selItem)
					m.selected--
				} else {
					end := len(m.list.Items()) - 1
					selItem := m.list.Items()[m.selected]
					for i := range m.list.Items() {
						if i != end {
							nextItem := m.list.Items()[i+1]
							m.list.SetItem(i, nextItem)
						} else {
							m.list.SetItem(end, selItem)
						}
					}
					m.selected, m.cursor = end, end
				}
			}
			m.list.CursorUp()
		case "down", "j":
			if m.selected != -1 {
				if m.cursor != len(m.list.Items())-1 {
					selItem := m.list.Items()[m.selected]
					swapItem := m.list.Items()[m.selected+1]
					m.list.SetItem(m.selected, swapItem)
					m.list.SetItem(m.selected+1, selItem)
					m.selected++
				} else {
					end := len(m.list.Items()) - 1
					lastItem := m.list.Items()[end]
					for i := end; i >= 0; i-- {
						if i != 0 {
							prevItem := m.list.Items()[i-1]
							m.list.SetItem(i, prevItem)
						} else {
							m.list.SetItem(0, lastItem)
						}
					}
					m.selected, m.cursor = 0, 0
				}
			}
			m.list.CursorDown()
		case "enter":
			if m.selected == -1 {
				m.selected = m.cursor
			} else {
				m.selected = -1
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "backspace":
			if m.selected == -1 {
				return m, tea.Cmd(func() tea.Msg { return exitDetail{} })
			}
			m.selected = -1
		case "left", "h":
			return m, tea.Cmd(func() tea.Msg { return unfocusDetail{} })
		}
		return m, tea.Cmd(func() tea.Msg { return handledCommand{} })
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *cascadeModel) View() string {
	return renderCascadeOrderView(m)
}

// type bookmarksModel struct {
// 	name  string
// 	list  list.Model
// 	focus bool
// }
//
// func (m *bookmarksModel) ShowFocus(focus bool) {
// 	m.focus = focus
// }
//
// func (m *bookmarksModel) Init() tea.Cmd {
// 	return nil
// }
//
// func (m *bookmarksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	// case tea.WindowSizeMsg:
// 	// 	m.width = msg.Width
// 	// 	m.height = msg.Height
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		case "esc", "backspace":
// 			// return resetSettings(), setSizeMsg(m.width, m.height)
// 		}
// 	}
// 	return m, nil
// }
//
// func (m *bookmarksModel) View() string {
// 	return renderManageBookmarksView(m)
// }
//
// type versionModel struct {
// 	name  string
// 	list  list.Model
// 	focus bool
// }
//
// func (m *versionModel) ShowFocus(focus bool) {
// 	m.focus = focus
// }
//
// func (m *versionModel) Init() tea.Cmd {
// 	return nil
// }
//
// func (m *versionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	// case tea.WindowSizeMsg:
// 	// 	m.width = msg.Width
// 	// 	m.height = msg.Height
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		case "esc", "backspace":
// 			// return resetSettings(), setSizeMsg(m.width, m.height)
// 		}
// 	}
// 	return m, nil
// }
//
// func (m *versionModel) View() string {
// 	return renderVersionView(m)
// }
