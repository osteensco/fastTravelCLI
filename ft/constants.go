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
