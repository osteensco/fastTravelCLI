package settings

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type navigateBackMsg struct{}

func back() tea.Cmd {
	return func() tea.Msg {
		return navigateBackMsg{}
	}
}

type settingsModel struct {
	docStyle lg.Style
	list          list.Model
	selectedView  ViewFunc
	selectedModel tea.Model
	width int
	height int
}

type item struct {
	name, desc string
	view       ViewFunc
	model      int
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.name }

var (
	modelList = []tea.Model{
		cascadeModel{
			docStyle: lg.NewStyle().Margin(1, 2),
			items: []string{"one", "two", "three", "four"}, //import from ft settings
			cursor: 0,
			selected: -1,
		},
		bookmarksModel{
			name: "BOOKMARKS",
		},
		versionModel{
			name: "VERSION INFO",
		},
	}

	settingsOptions = []list.Item{

		item{
			name: "Cascade Order", 
			desc: "Determine the ordering in which fastTravelCLI resolves a query.", 
			view: renderCascadeOrderView, 
			model: 0,
		},

		item{
			name: "Manage Bookmarks", 
			desc: "View and manage bookmarks saved with fastTravelCLI.", 
			view: renderManageBookmarksView, 
			model: 1,
		},

		item{
			name: "Version", 
			desc: "View fastTravelCLI version information.", 
			view: renderVersionView, 
			model: 2,
		},

	}

	settingsList = list.New(settingsOptions, list.NewDefaultDelegate(), 0, 0)
)

func (m settingsModel) Init() tea.Cmd {
	return nil
}

func (m settingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := m.docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		// exit the program
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		// select a setting
		case "enter":
			i := m.list.SelectedItem().(item)
			m.selectedView = i.view

			m.selectedModel = modelList[i.model]
			return m.selectedModel, setSizeMsg(m.width, m.height)
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// The view is just one big string that gets passed to the UI for rendering
func (m settingsModel) View() string {
	if m.selectedView != nil {
		return m.selectedView(m.selectedModel)
	}
	return renderSettingsView(m)
}






type cascadeModel struct {
	docStyle lg.Style
	items []string
	cursor int
	selected int
	width        int
	height       int
}

func (m cascadeModel) Init() tea.Cmd {
	return nil
}

func (m cascadeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected == -1 {
				if m.cursor != 0 {
					m.cursor = m.cursor-1
				} else {
					m.cursor = len(m.items)-1
				}
			} else {
				if m.cursor != 0 {
					m.items[m.selected], m.items[m.selected-1] = m.items[m.selected-1], m.items[m.selected]
					m.selected--
					m.cursor = m.selected
				} else {
					end := len(m.items)-1
					selected := m.items[m.selected]
					for i, _ := range m.items {
						if i != end {
							m.items[i] = m.items[i+1]
						} else {
							m.items[end] = selected
						}
					}
					m.selected, m.cursor = end, end
				}
			}
		case "down", "j":
			if m.selected == -1 {
				if m.cursor != len(m.items)-1 {
					m.cursor = m.cursor+1
				} else {
					m.cursor = 0
				}
			} else {
				if m.cursor != len(m.items)-1 {
					m.items[m.selected], m.items[m.selected+1] = m.items[m.selected+1], m.items[m.selected]
					m.selected++
					m.cursor = m.selected
				} else {
					end := len(m.items)-1
					lastItem := m.items[end]
					for i := end; i >= 0; i-- {
						if i != 0 {
							m.items[i] = m.items[i-1]
						} else {
							m.items[0] = lastItem
						}
					}
					m.selected, m.cursor = 0, 0
				}
			}
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
				return resetSettings(), setSizeMsg(m.width, m.height)
			}
			m.selected = -1
		}
	}
	return m, nil
}

func (m cascadeModel) View() string {
	return renderCascadeOrderView(m)
}






type bookmarksModel struct {
	name         string
	width        int
	height       int
}

func (m bookmarksModel) Init() tea.Cmd {
	return nil
}

func (m bookmarksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "backspace":
			return resetSettings(), setSizeMsg(m.width, m.height)
		}
	}
	return m, nil
}

func (m bookmarksModel) View() string {
	return renderManageBookmarksView(m)
}






type versionModel struct {
	name         string
	width        int
	height       int
}

func (m versionModel) Init() tea.Cmd {
	return nil
}

func (m versionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc", "backspace":
			return resetSettings(), setSizeMsg(m.width, m.height)
		}
	}
	return m, nil
}

func (m versionModel) View() string {
	return renderVersionView(m)
}
