package data

import (
	"encoding/json"
	"fmt"
	"os"
)

// TODO: Settings to add
// 	- toggle for history across active sessions (think tmux workflows)
//

type Settings struct {
	QueryOrder   []string `json:"query_order"`
	UpdatePrompt bool     `json:"update_prompt"`
}

func GenerateDefaultSettings() *Settings {
	return &Settings{
		QueryOrder:   []string{"bookmark", "CDPATH", "relative"},
		UpdatePrompt: false,
	}
}

func ReadSettings(settingsFile *os.File) (*Settings, error) {
	// use default settings if we don't have any settings
	settings := GenerateDefaultSettings()

	info, err := settingsFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("Error getting file info: %w", err)
	}

	// only read settings in if the settings file is populated
	if info.Size() != 0 {
		err = json.NewDecoder(settingsFile).Decode(settings)
		if err != nil {
			return nil, err
		}
	}

	return settings, nil
}

func WriteSettings(settingsFile *os.File, settings *Settings) error {
	// err := settingsFile.Truncate(0)
	// if err != nil {
	// 	return fmt.Errorf("Error truncating file: %w", err)
	// }
	//
	// _, err = settingsFile.Seek(0, 0)
	// if err != nil {
	// 	return fmt.Errorf("Error seeking to beginning of file: %w", err)
	// }
	//
	// err = json.NewEncoder(settingsFile).Encode(settings)
	// if err != nil {
	// 	return err
	// }
	
	fmt.Println("Write Success!")
	fmt.Println(settings.QueryOrder)

	return nil
}
