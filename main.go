package main

import (
	"fmt"
	"os"
	"path/filepath"
    ft "github.com/osteensco/fastTravelCLI/src"
)

// TODO
// ***FIX TESTS TO ACCOUNT FOR EDGE CASES***
// - tests
//      - currently focuses on happy path
//      - add tests for edge cases

// - flags
//      - flags for commands, like -h
//      - global flags like -v
//          - essentially shortcuts for other commands

// - features:
//      - ft version
//      - ft update
//      - ft set ?
//          - open live fuzzy finder
//          - interactively set a key
//          - results matching current entry can be selected with arrow keys and pressing enter
//      - ft swap [key1] [key2]
//          - swaps the dirs of the two keys given
//      - ft ?
//          - ask fastTravel if the curr dir is saved
//      - remember prev n directories? n=10?15?20?
//      - add a spinner?






func main() {

	
    // read in bin file
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	dataDirPath := filepath.Dir(exePath)
	dataPath := dataDirPath + "\\fastTravel.bin"
    file := ensureData(dataPath)
    defer file.Close()
    allPaths := readMap(file)

	// sanitize input
    inputCommand, err := passCmd(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	action := inputCommand[0]

	// execute user provided action
	exeCmd, ok := availCmds[action]
	if !ok {
		fmt.Println("Invalid command, use 'help' for available commands.")
		os.Exit(1)
	}

	data := cmdArgs{
		cmd:      inputCommand,
		allPaths: allPaths,
        file: file,
        rdr: os.Stdin,
	}
    
	exeCmd(data)

}
