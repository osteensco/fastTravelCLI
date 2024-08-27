package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteensco/fastTravelCLI/ft"
)

// TODO
// ***FIX TESTS TO ACCOUNT FOR EDGE CASES***
// - tests
//      - currently focuses on happy path
//      - add tests for edge cases

// - new features:
//      - lua bindings to change command syntax
//          - ie user wants to change "-set" to "-s"
//      - track directories visited
//          - use dirs -v, pushd, popd
//              - ft seems to already add to session stack
//              - currently, when popd is called, path is removed from stack
//                  - add global stack variable to config file to capture popd output
//                  - ft bash function will need to implement push/pop stack operations
//          - ft < and ft > to navigate stack
//      - ft some/path/given should simply pass path to stdout (mimick cd)
//          - ft ~[key]/subdir should evaluate to [keyValue]/subdir and pass to stdout
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
		fmt.Println("Invalid command, use 'help' for available commands.")
		os.Exit(1)
	}

	data := ft.NewCmdArgs(inputCommand, allPaths, file, os.Stdin)

	exeCmd(data)

}
