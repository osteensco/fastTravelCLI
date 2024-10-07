package ft

import (
	"io"
	"os"
)

// ft command api
type CmdArgs struct {
	cmd      []string
	allPaths map[string]string
	file     *os.File
	rdr      io.Reader
}

func NewCmdArgs(inputCmd []string, allPaths map[string]string, file *os.File, rdr io.Reader) *CmdArgs {
	return &CmdArgs{inputCmd, allPaths, file, rdr}
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
}

var CmdDesc = map[string]string{
	"key":      "change directory to provided key's path - Usage: ft [key]",
	"-set":     "set current directory path to provided key - Usage: ft -set [key]",
	"-ls":      "display all current key value pairs - Usage: ft -ls",
	"-rm":      "deletes provided key - Usage: ft -rm [key]",
	"-rn":      "renames key to new key - Usage: ft -rn [key] [new key]",
	"]":        "navigate history forwards - Usage: ft ]",
	"[":        "navigate history backwards - Usage: ft [",
	"-is":      "know the directory variable if set for a directory",
	"-help":    "you are here :) - Usage: ft -help, -h",
	"-version": "print current version of fastTravelCLI - Usage: ft -version, -v",
}

// default value for version
var Version string = "development"

// default value for logo
var Logo string = "fastTravelCLI"

// responses
const (
	noLocationsSetMsg         = "No fast travel locations set, set locations by navigating to desired destination directory and using 'ft -set <key>' \n"
	invalidDirectoryMsg       = "Provided path '%s' evaluates to '%s' which is not a valid directory. Use 'ft -ls' to see all saved destinations. \n"
	unrecognizedKeyMsg        = "Did not recognize key '%s', use 'ft -ls' to see all saved destinations. If this is a relative path use './%s' or '%s/'. \n"
	pathAlreadyExistsMsg      = "Path '%s' already exists with key '%s', overwrite key '%s' \n"
	abortedOverwriteKeyMsg    = "Aborted overwriting of key '%s' \n"
	renamedKeyMsg             = "Renamed key '%s' to '%s' whose value is '%s' \n"
	keyAlreadyExistsMsg       = "Key '%s' already exists with value '%s', overwrite key '%s'?(y/n) "
	addKeyMsg                 = "Added destination '%s': '%s' \n"
	keyDoesNotExistMsg        = "Key '%s' does not exist. Run 'ft -ls' to see all keys. \n"
	verifyRemoveMsg           = "Are you sure you want to remove the key '%s'? (y/n) "
	abortRemoveKeyMsg         = "Aborted removal of key '%s' \n"
	removeKeyMsg              = "Removed '%s' destination \n"
	renameKeyAlreadyExistsMsg = "Key '%s' already exists, please choose something else. \n"
	renameKeyDoesNotExistMsg  = "Cannot rename '%s', key does not exist. Run 'ft -ls' to see all keys. \n"
	verifyRenameMsg           = "Are you sure you want to rename the key '%s' to '%s'? (y/n) "
	abortRenameKeyMsg         = "Aborted renaming of key '%s' to '%s'. \n"
)

//
