package ft

import (
	"io"
	"os"
)

// ft api
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
var AvailCmds = map[string]func(data *CmdArgs){
	"_":     changeDirectory,
	"-set":  setDirectoryVar,
	"-ls":   displayAllPaths,
	"-rm":   removeKey,
	"-rn":   renameKey,
	"-help": showHelp,
	"-h":    showHelp,
}

var CmdDesc = map[string]string{
	"[key]": "change directory to provided key's path - Usage: ft [key]",
	"-set":  "set current directory path to provided key - Usage: ft -set [key]",
	"-ls":   "display all current key value pairs - Usage: ft -ls",
	"-rm":   "deletes provided key - Usage: ft -rm [key]",
	"-rn":   "renames key to new key - Usage: ft -rn [key] [new key]",
	"-help": "you are here :) - Usage: ft -help",
}
