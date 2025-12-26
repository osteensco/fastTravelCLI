package tui

import (
	ftdata "github.com/osteensco/fastTravelCLI/data"
	"github.com/osteensco/fastTravelCLI/tui/settings"
)

type settingsTUI struct {
	Settings ftdata.Settings
}

func (s *settingsTUI) Run() error {
	err := settings.Run(&s.Settings)
	if err != nil {
		return err
	}
	return nil
}

type tui struct {
	Settings settingsTUI
}

func Init() tui {
	return tui{
		settingsTUI{},
	}
}
