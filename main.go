package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/osteensco/fastTravelCLI/ft"
)

// fastTravelCLI main process
func main() {

	// identify exe path to establish a working directory and find dependency files
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// handle piped args
	err = ft.PipeArgs(&os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// sanitize user input
	inputCommand, err := ft.PassCmd(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	action := inputCommand.Cmd

	// grab command from registry
	cmd, ok := ft.AvailCmds[action]
	if !ok {
		fmt.Printf("Invalid command '%s', use 'ft -h' for available commands. \n", action)
		return
	}

	var dataDirPath string
	var dataPath string
	var file *os.File
	var allPaths map[string]string

	if cmd.LoadData {
		// find persisted keys or create file to persist keys
		dataDirPath = filepath.Dir(exePath)
		dataPath = fmt.Sprintf("%s/fastTravel.bin", dataDirPath)

		file, err = ft.EnsureData(dataPath)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer file.Close()

		// read keys into memory
		allPaths, err = ft.ReadMap(file)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	// manifest API
	data := ft.NewCmdAPI(dataDirPath, inputCommand, allPaths, file, os.Stdin)

	// execute user provided action
	err = cmd.Callback(data)
	if err != nil {
		fmt.Println("fastTravelCLI returned an error: ", err)
		return
	}

}
