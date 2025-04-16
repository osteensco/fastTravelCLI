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
	H bool
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
const HelpLineStrFormat = "  %-18s %s\n"

const HelpDescExamples = `
    Examples: 
	  ft projects                 → cd into the directory saved as 'projects'
	  ft -set docs=~/Documents    → set key 'docs' to '~/Documents'
	  ft -rn -y old new           → rename 'old' key to 'new' and skip confirmation prompt
    `

// TODO
var HelpUsageMappings = map[string]string{
	"key":     "ft <key>",
	"set":     "ft -set [-y] <key>=[path]",
	"ls":      "ft ls",
	"rm":      "ft -rm [-y] <key>",
	"rn":      "ft -rn [-y] <key> <new key>",
	"edit":    "ft -edit [-y] <path> <new dir name>",
	"]":       "ft ]",
	"[":       "ft [",
	"hist":    "ft -hist",
	"is":      "ft -is",
	"version": "ft -version, -v",
	"update":  "ft -update, -u [version][nightly]",
	"help":    "ft -help, -h [command]",
}

var CmdDesc = []map[string]string{
	{
		HelpUsageMappings["key"]: "Change directory to a saved key",
	},
	{
		HelpUsageMappings["set"]: "Save a key to a path (defaults to CWD if no path)",
	},
	{
		HelpUsageMappings["ls"]: "List all saved key-path pairs",
	},
	{
		HelpUsageMappings["rm"]: "Remove a saved key",
	},
	{
		HelpUsageMappings["rn"]: "Renames an existing key",
	},
	{
		HelpUsageMappings["edit"]: "Updates a given directory to it's new name for any key assigned to it or a child directory",
	},
	{
		HelpUsageMappings["]"]: "Navigate history forwards",
	},
	{
		HelpUsageMappings["["]: "Navigate history backwards",
	},
	{
		HelpUsageMappings["hist"]: "Show directory history with fzf",
	},
	{
		HelpUsageMappings["is"]: "Show the key associated with CWD",
	},
	{
		HelpUsageMappings["version"]: "Show current version",
	},
	{
		HelpUsageMappings["update"]: "Update fastTravelCLI, optionally specify version or nightly, defaults to latest",
	},
	{
		HelpUsageMappings["help"]: "Show this help message",
	},
}

// Detailed help docs
var DetailedCmdDescriptions = [][]string{
	{HelpUsageMappings["key"],
		`Description:
  Changes the current directory to the one associated with the given key.

Accepted Formats for <key>:
  - A valid saved key
  - A key/subdir format (e.g., pr/docs/work)
  - A relative path (e.g., ./mydir)
  - An absolute path (e.g., /home/user/mydir)

Examples:
  ft projects         # cd into the path saved as "projects"
  ft dev/mytool       # cd into a subdirectory under the "dev" key`},

	{HelpUsageMappings["set"],
		`Description:
  Saves a key to a directory path.
  If no path is specified, the current working directory is used.

Arguments:
  key         The name to assign the path to.
  path        (Optional) The path to associate with the key.

Flags:
  -y    Skip confirmation prompt.

Examples:
  ft -set docs                        # Set "docs" to the current working directory
  ft -set code=~/Projects             # Set "code" to a specific path
  ft -set os=~/Projects/opensource    # Set "os" to specified path without confirmation prompt`},
	{HelpUsageMappings["ls"],
		`Description:
  Lists all saved key-path pairs.`},
	{HelpUsageMappings["rm"],
		`Description:
  Deletes a saved key.

Flags:
  -y    Skip confirmation prompt.

Examples:
  ft -rm archive       # Prompts for confirmation
  ft -rm archive -y    # Removes without confirmation`},
	{HelpUsageMappings["rn"],
		`Description:
  Renames an existing key to a new one.

Flags:
  -y    Skip confirmation prompt.

Examples:
  ft -rn old new       # "old" is renamed to "new"
  ft -rn old new -y    # "old" is renamed to "new" without confirmation`},
	{HelpUsageMappings["edit"],
		`Description:
  Renames a directory (and all child directories) associated with a key or path.
  Updates any keys pointing to subdirectories of the renamed folder.

Accepted Formats for <old>:
  - A relative path (e.g., ./mydir)
  - An absolute path (e.g., /home/user/mydir)
  - A key/subdir format (e.g., projects/docs)

Examples:
  ft -edit ./api api-v2
  ft -edit projects/docs documents`},
	{HelpUsageMappings["]"],
		`Description:
  Navigate forward in directory history.`},
	{HelpUsageMappings["["],
		`Description:
  Navigate backward in directory history.`},
	{HelpUsageMappings["hist"],
		`Description:
  Shows an interactive history for fuzzy selection.

Notes:
  Requires 'fzf' and 'tree' to be installed.`},
	{HelpUsageMappings["version"],
		`Description:
  Displays the current version of fastTravelCLI.`},
	{HelpUsageMappings["is"],
		`Description:
  Displays the key associated with the current working directory, if one exists.`},
	{HelpUsageMappings["update"],
		`Description:
  Updates fastTravelCLI to the latest version.
  If a version is specified, attempts to install that version.
  Use "nightly" to install the latest pre-release build.

Examples:
  ft -update
  ft -update v.0.2.1
  ft -update nightly`},
	{"This is a standard directory navigation command.", "ft can replace cd entirely as it inherits cd's commands."},
}

const HelpCmdMsg = "Use ft -help information on other commands"

var DetailedCmdDescMapping = map[string]string{
	"_":        CreateCmdHelpDoc(DetailedCmdDescriptions[0]),
	"-set":     CreateCmdHelpDoc(DetailedCmdDescriptions[1]),
	"-ls":      CreateCmdHelpDoc(DetailedCmdDescriptions[2]),
	"-rm":      CreateCmdHelpDoc(DetailedCmdDescriptions[3]),
	"-rn":      CreateCmdHelpDoc(DetailedCmdDescriptions[4]),
	"-edit":    CreateCmdHelpDoc(DetailedCmdDescriptions[5]),
	"-]":       CreateCmdHelpDoc(DetailedCmdDescriptions[6]),
	"-[":       CreateCmdHelpDoc(DetailedCmdDescriptions[7]),
	"-hist":    CreateCmdHelpDoc(DetailedCmdDescriptions[8]),
	"-version": CreateCmdHelpDoc(DetailedCmdDescriptions[9]),
	"-v":       CreateCmdHelpDoc(DetailedCmdDescriptions[9]),
	"-is":      CreateCmdHelpDoc(DetailedCmdDescriptions[10]),
	"-update":  CreateCmdHelpDoc(DetailedCmdDescriptions[11]),
	"-u":       CreateCmdHelpDoc(DetailedCmdDescriptions[11]),
	"--":       CreateCmdHelpDoc(DetailedCmdDescriptions[12]),
	"-..":      CreateCmdHelpDoc(DetailedCmdDescriptions[12]),
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
