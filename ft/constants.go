package ft

import (
    "os"
    "io"
)



// api
type cmdArgs struct {
	cmd      []string
	allPaths map[string]string
    file *os.File
    rdr io.Reader
}

// map of available commands
var availCmds = map[string]func(data cmdArgs){
	"to":  changeDirectory,
	"set": setDirectoryVar,
	"ls":  displayAllPaths,
    "rm": removeKey,
    "rn": renameKey,
    "help": showHelp,
}

var cmdDesc = map[string]string {
    "to": "change directory to provided key's path - Usage: ft to [key]",
    "set": "set current directory path to provided key - Usage: ft set [key]",
    "ls": "display all current key value pairs - Usage: ft ls",
    "rm": "deletes provided key - Usage: ft rm [key]",
    "rn": "renames key to new key - Usage: ft rn [key] [new key]",
    "help": "you are here :) - Usage: ft help",
}
