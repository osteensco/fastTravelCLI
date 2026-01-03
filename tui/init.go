package tui

import (
	"os"

	ftdata "github.com/osteensco/fastTravelCLI/data"
	"github.com/osteensco/fastTravelCLI/tui/settings"
)

type settingsTUI struct {
	Settings *ftdata.Settings
	File     *os.File
}

func (s *settingsTUI) Run() error {
	err := settings.Run(s.Settings, s.File)
	if err != nil {
		return err
	}
	return nil
}

type tui struct {
	Settings settingsTUI
}

func Init(settings *ftdata.Settings, file *os.File) tui {
	return tui{
		settingsTUI{Settings: settings, File: file},
	}
}
