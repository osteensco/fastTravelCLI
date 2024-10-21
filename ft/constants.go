package ft

import (
	"io"
	"os"
)

// ft command api
type CmdArgs struct {
	wkDir    string
	cmd      []string
	allPaths map[string]string
	file     *os.File
	rdr      io.Reader
}

func NewCmdArgs(ftDir string, inputCmd []string, allPaths map[string]string, file *os.File, rdr io.Reader) *CmdArgs {
	return &CmdArgs{ftDir, inputCmd, allPaths, file, rdr}
}

// map of available commands
var AvailCmds = map[string]func(data *CmdArgs) error{
	"_":        changeDirectory,
	"-set":     setDirectoryVar,
	"-ls":      displayAllPaths,
	"-rm":      removeKey,
	"-rn":      renameKey,
	"-help":    showHelp,
	"-h":       showHelp,
	"-]":       passToShell,
	"-[":       passToShell,
	"-..":      passToShell,
	"--":       passToShell,
	"-version": showVersion,
	"-v":       showVersion,
	"-is":      showDirectoryVar,
	"-update":  updateFT,
	"-u":       updateFT,
}

var CmdDesc = map[string]string{
	"key":      "change directory to provided key's path - Usage: ft [key]",
	"-set":     "set current directory path to provided key - Usage: ft -set [key]",
	"-ls":      "display all current key value pairs - Usage: ft -ls",
	"-rm":      "deletes provided key - Usage: ft -rm [key]",
	"-rn":      "renames key to new key - Usage: ft -rn [key] [new key]",
	"]":        "navigate history forwards - Usage: ft ]",
	"[":        "navigate history backwards - Usage: ft [",
	"-is":      "identify the key that the current working directory is saved to if it is saved to a key - Usage: ft -is",
	"-help":    "you are here :) - Usage: ft -help, -h",
	"-version": "print current version of fastTravelCLI - Usage: ft -version, -v",
	"-update":  "update fastTravel - Usage: ft -update, -u [version](optional)",
}

// default value for version
var Version string = "development"

// default value for logo
var Logo string = "fastTravelCLI"

// user messages
const (
	NoLocationsSetMsg         = "No fast travel locations set, set locations by navigating to desired destination directory and using 'ft -set <key>' \n"
	InvalidDirectoryMsg       = "Provided path '%s' evaluates to '%s' which is not a valid directory. Use 'ft -ls' to see all saved destinations. \n"
	UnrecognizedKeyMsg        = "Did not recognize key '%s', use 'ft -ls' to see all saved destinations. If this is a relative path use './%s' or '%s/'. \n"
	PathAlreadyExistsMsg      = "Path '%s' already exists with key '%s', overwrite key '%s' \n"
	AbortedOverwriteKeyMsg    = "Aborted overwriting of key '%s' \n"
	RenamedKeyMsg             = "Renamed key '%s' to '%s' whose value is '%s' \n"
	KeyAlreadyExistsMsg       = "Key '%s' already exists with value '%s', overwrite key '%s'?(y/n) "
	AddKeyMsg                 = "Added destination '%s': '%s' \n"
	KeyDoesNotExistMsg        = "Key '%s' does not exist. Run 'ft -ls' to see all keys. \n"
	VerifyRemoveMsg           = "Are you sure you want to remove the key '%s'? (y/n) "
	AbortRemoveKeyMsg         = "Aborted removal of key '%s' \n"
	RemoveKeyMsg              = "Removed '%s' destination \n"
	RenameKeyAlreadyExistsMsg = "Key '%s' already exists, please choose something else. \n"
	RenameKeyDoesNotExistMsg  = "Cannot rename '%s', key does not exist. Run 'ft -ls' to see all keys. \n"
	VerifyRenameMsg           = "Are you sure you want to rename the key '%s' to '%s'? (y/n) "
	AbortRenameKeyMsg         = "Aborted renaming of key '%s' to '%s'. \n"
	IsKeyMsg                  = "Directory %s is saved to key : %s \n"
	IsNotKeyMsg               = "No key was found for the specified path: %s \n"
)
