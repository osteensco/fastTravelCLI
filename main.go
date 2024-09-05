package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteensco/fastTravelCLI/ft"
)

// TODO

// - tests
//      - add more tests for edge cases
//      - add shell ft function tests

// - fixes
//      - prettier console prints
//      - currently cd will not affect upperStack. should this be the same using ft?

// - new features:
//      - ft .. (mimick cd ..)
//      - ft - (mimick cd -)
//          - should this affect upperStack?
//      - ft -version, -v
//          - have this also disply ascii art
//      - ft -update, -u
//      - ft -set ?
//          - open live fuzzy finder
//          - interactively set a key
//          - results matching current entry can be selected with arrow keys and pressing enter
//      - ft ?
//          - ask fastTravel if the curr dir is saved

func main() {

	// read in bin file
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	dataDirPath := filepath.Dir(exePath)
	dataPath := fmt.Sprintf("%s/fastTravel.bin", dataDirPath)
	file := ft.EnsureData(dataPath)
	defer file.Close()
	allPaths := ft.ReadMap(file)

	// sanitize input
	inputCommand, err := ft.PassCmd(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	action := inputCommand[0]

	// execute user provided action
	exeCmd, ok := ft.AvailCmds[action]
	if !ok {
		fmt.Printf("Invalid command '%s', use 'help' for available commands.\n", action)
		os.Exit(1)
	}

	data := ft.NewCmdArgs(inputCommand, allPaths, file, os.Stdin)

	exeCmd(data)

}
