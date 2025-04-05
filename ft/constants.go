package ft

import (
	"io"
	"os"
)

// ft command api
type CmdAPI struct {
	wkDir    string
	cmd      *Cmd
	allPaths map[string]string
	file     *os.File
	rdr      io.Reader
}

func NewCmdAPI(ftDir string, inputCmd *Cmd, allPaths map[string]string, file *os.File, rdr io.Reader) *CmdAPI {
	return &CmdAPI{ftDir, inputCmd, allPaths, file, rdr}
}

// struct used to identify flags that were provided with a given command
type CmdFlags struct {
	y bool
}

// struct used to dissect and organize a command into it's individual components
type Cmd struct {
	Flags CmdFlags
	Cmd   string
	Args  []string
}

func NewCmd(args *[]string) *Cmd {
	return &Cmd{
		// flags and cmd will be empty defaults
		// args is explicit so that enough space is allocated to unerlying array for slight optimization
		Args: make([]string, 0, len(*args)),
	}
}

// map of available commands
var AvailCmds = map[string]struct {
	Callback func(data *CmdAPI) error
	LoadData bool
}{
	"_":        {changeDirectory, true},
	"-set":     {setDirectoryVar, true},
	"-setf":    {setDirectoryVar, true}, // TODO delete me
	"-ls":      {displayAllPaths, true},
	"-rm":      {removeKey, true},
	"-rn":      {renameKey, true},
	"-edit":    {editPath, true},
	"-help":    {showHelp, false},
	"-h":       {showHelp, false},
	"-]":       {passToShell, false},
	"-[":       {passToShell, false},
	"-hist":    {passToShell, false},
	"-..":      {passToShell, false},
	"--":       {passToShell, false},
	"-version": {showVersion, false},
	"-v":       {showVersion, false},
	"-is":      {showDirectoryVar, true},
	"-update":  {updateFT, false},
	"-u":       {updateFT, false},
}

var CmdDesc = map[string]string{
	"key":      "change directory to provided key's path - Usage: ft [key]",
	"-set":     "set key to a directory path, if no directory path is given attempts to set key to CWD - Usage: ft -set [key], ft -set [key]=[path]",
	"-setf":    "force set key to a directory path, if no directory path is given force set key to CWD - Usage: ft -setf [key], ft -setf [key]=[path]",
	"-ls":      "display all current key value pairs - Usage: ft -ls",
	"-rm":      "deletes provided key - Usage: ft -rm [key]",
	"-rn":      "renames key to new key - Usage: ft -rn [key] [new key]",
	"-edit":    "renames a given directory and all child directories assigned to a key to a new folder name - Usage: ft -edit mydir mynewdir",
	"]":        "navigate history forwards - Usage: ft ]",
	"[":        "navigate history backwards - Usage: ft [",
	"-hist":    "show history, provides fzf selection for navigation - Usage: ft -hist",
	"-is":      "identify the key that the current working directory is saved to if it is saved to a key - Usage: ft -is",
	"-help":    "you are here :) - Usage: ft -help, -h",
	"-version": "print current version of fastTravelCLI - Usage: ft -version, -v",
	"-update":  "update fastTravel - Usage: ft -update, -u [version](optional)",
}

// default value for version
var Version string = "development"

// default value for logo
var Logo string = "fastTravelCLI"

// GitHub endpoints for update command
var EndpointGH string = "https://api.github.com/repos/osteensco/fastTravelCLI/releases/tags/%s"
var EndpointLatestGH string = "https://api.github.com/repos/osteensco/fastTravelCLI/releases/latest"

// CLI command and resulting directory
var GitCloneCMD []string = []string{"git", "clone", "--branch", "", "https://github.com/osteensco/fastTravelCLI.git"}
var GitCloneDir string = "fastTravelCLI"

// Switch for testing
var UPDATEMOCK bool = false

// user messages
const (
	InvalidDirectoryMsg       = "Provided path '%s' evaluates to '%s' which is not a valid directory. Use 'ft -ls' to see all saved destinations. \n"
	UnrecognizedKeyMsg        = "Did not recognize key or relative path '%s', use 'ft -ls' to see all saved destinations. \n"
	PathAlreadyExistsMsg      = "Path '%s' already exists with key '%s', overwrite key '%s'? (y/n) \n"
	AbortedOverwriteKeyMsg    = "Aborted overwriting of key '%s'. \n"
	PathOverwriteMsg          = "The value of key '%s' has been overwritten and is now '%s'. \n"
	RenamedKeyMsg             = "Renamed key '%s' to '%s' whose value is '%s'. \n"
	KeyAlreadyExistsMsg       = "Key '%s' already exists with value '%s', overwrite key '%s'? (y/n) \n"
	AddKeyMsg                 = "Added destination '%s': '%s'. \n"
	KeyDoesNotExistMsg        = "Key '%s' does not exist. Run 'ft -ls' to see all keys. \n"
	VerifyRemoveMsg           = "Are you sure you want to remove the key '%s'? (y/n) \n"
	AbortRemoveKeyMsg         = "Aborted removal of key '%s'. \n"
	RemoveKeyMsg              = "Removed '%s' destination. \n"
	RenameKeyAlreadyExistsMsg = "Key '%s' already exists, please choose something else. \n"
	RenameKeyDoesNotExistMsg  = "Cannot rename '%s', key does not exist. Run 'ft -ls' to see all keys. \n"
	VerifyRenameMsg           = "Are you sure you want to rename the key '%s' to '%s'? (y/n) \n"
	AbortRenameKeyMsg         = "Aborted renaming of key '%s' to '%s'. \n"
	IsKeyMsg                  = "Directory %s is saved to key : %s. \n"
	IsNotKeyMsg               = "No key was found for the specified path: %s. \n"
)
