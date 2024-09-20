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
//      - cd does not affect upperStack. should this be the same using ft?

// - new features:

//      - ft -version, -v
//          - have this also disply ascii art
//      - ft -update, -u
//          - automagically pull and source if new version available
//      - ft -set ?
//          - open live fuzzy finder
//          - interactively set a key
//          - results matching current entry can be selected with arrow keys and pressing enter
//      - ft ?
//          - ask fastTravel if the curr dir is saved
//      - ft -load
//          - add key value pairs via json file
//          - useful for preconfigured setups like as part of a dockerfile
//      - ft -massedit [substring] [new substring]
//          - allows users to easily update directory names they changed
//          - consider an example, a directory called "programming stuff"
//              - I set a key to "programming stuff/coolProject"
//              - I realize having spaces in a dir name is silly so I rename it to programmingStuff
//              - this key is now broken in ft, because ft would have no way of knowing it changed
//              - it's easy to reset one key, but if I have 10 keys saved to sub dirs of something I changed that sucks
//              - this command would address this problem graciously
//      - ft -q [some query]
//          - pass a query and return keys and values that match
//      - ft -hist
//          - display history stack
//          - include simple indicator of where user is currently in history stack

func main() {

	// read in bin file
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	dataDirPath := filepath.Dir(exePath)
	dataPath := fmt.Sprintf("%s/fastTravel.bin", dataDirPath)

	file, err := ft.EnsureData(dataPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()

	allPaths, err := ft.ReadMap(file)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// sanitize input
	inputCommand, err := ft.PassCmd(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	action := inputCommand[0]

	// execute user provided action
	exeCmd, ok := ft.AvailCmds[action]
	if !ok {
		fmt.Printf("Invalid command '%s', use 'ft -h' for available commands. \n", action)
		return
	}

	data := ft.NewCmdArgs(inputCommand, allPaths, file, os.Stdin)

	err = exeCmd(data)
	if err != nil {
		fmt.Println("fastTravel returned an error: ", err)
		return
	}

}
