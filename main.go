package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteensco/fastTravelCLI/ft"
)

// fastTravel main process
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
		fmt.Println("fastTravelCLI returned an error: ", err)
		return
	}

}
