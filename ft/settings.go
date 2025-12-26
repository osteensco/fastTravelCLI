package ft

// TODO create settings data structures

// Settings

// - toggle for history across active sessions (think tmux workflows)
// - auto updates (or check for updates)




type Settings struct {
	cascadeOrder []string
}

func NewSettings() *Settings {
	return &Settings{}
}
