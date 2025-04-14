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
	Y bool
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
		// args is explicit so that enough space is allocated to underlying array for slight optimization
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
	"-ls":      {displayAllPaths, true},
	"-rm":      {removeKey, true},
	"-rn":      {renameKey, true},
	"-edit":    {editPath, true},
	"-help":    {showHelp, false},
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

// Help docs
var HelpLineStrFormat = "  %-18s %s\n"

// TODO
var HelpDescExamples = ""

// TODO
var HelpUsageMappings = map[string]string{
	"set": "ft -set [-y] <key>=[path]",
}

var CmdDesc = []map[string]string{
	{
		"ft <key>": "Change directory to a saved key",
	},
	{
		HelpUsageMappings["set"]: "Save a key to a path (defaults to CWD if no path)",
	},
	{
		"ft -ls": "List all saved key-path pairs",
	},
	{
		"ft -rm [-y] <key>": "Remove a saved key",
	},
	{
		"ft -rn [-y] <key> <new key>": "Renames an existing key",
	},
	{
		"ft -edit [-y] <path> <new dir name>": "Updates a given directory to it's new name for any key assigned to it or a child directory",
	},
	{
		"ft ]": "Navigate history forwards",
	},
	{
		"ft [": "Navigate history backwards",
	},
	{
		"ft -hist": "Show directory history with fzf",
	},
	{
		"ft -is": "Show the key associated with CWD",
	},
	{
		"ft -version, -v": "Show current version",
	},
	{
		"ft -update, -u [version][nightly]": "Update fastTravelCLI, optionally specify version or nightly, defaults to latest",
	},
	{
		"ft -help": "Show this help message",
	},
}

// Detailed help docs

// TODO
var DetailedCmdDescriptions = [][]string{
	{HelpUsageMappings["set"], "A detailed description here"},
}

var DetailedCmdDescMapping = map[string]string{
	"-set":     "",
	"-ls":      "",
	"-rm":      "",
	"-rn":      "",
	"-edit":    "",
	"-help":    "",
	"-]":       "",
	"-[":       "",
	"-hist":    "",
	"-..":      "",
	"--":       "",
	"-version": "",
	"-v":       "",
	"-is":      "",
	"-update":  "",
	"-u":       "",
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
	VerifyEditMsg             = "Are you sure you want to change '%s: %s' to '%s: %s'? (y/n) \n"
	AbortEditMsg              = "Aborted replacing '%s: %s' with '%s: %s'. \n"
	PathIsNotValidDirWarn     = "Warning: Path %s is not a valid directory. \n"
)
