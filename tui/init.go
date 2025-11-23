package tui

import "github.com/osteensco/fastTravelCLI/tui/settings"

type settingsTUI struct{}

func (s *settingsTUI) Run() error {
	err := settings.Run()
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
