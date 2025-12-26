package data

// TODO create settings data structures

// Settings

// - toggle for history across active sessions (think tmux workflows)
// - auto updates (or check for updates)

type Settings struct {
	QueryOrder []string
}


func NewSettings() *Settings {
	return &Settings{}
}

func ReadInSettings(settings *Settings) {

}

func WriteSettings(settings *Settings) {

}
