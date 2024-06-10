package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func passCmd(args []string) ([]string, error) {
	
    if args[1] != "ls" && len(args) <= 2 {
		return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ftrav <command> <path/key>", args))
	}
	return args[1:], nil

}

func readMap(jsonPath string) map[string]string {
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer file.Close()

	var pathMap map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pathMap); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return pathMap

}

func changeDirectory(data cmdArgs) {

	if len(data.allPaths) == 0 {
		fmt.Printf("No fast travel locations set, set locations by navigating to desired destination directory and using 'ftrav set <key>' ")
		os.Exit(1)
	}

	path := data.allPaths[data.cmd[1]]
	distro := os.Getenv("WSL_DISTRO_NAME")
	if len(distro) == 0 {

		prefix := "/mnt/"
		path = strings.Replace(path, ":", "", 1)
		path = strings.Replace(path, "\\", "/", -1)
		path = strings.ToLower(prefix + path)

	}

    fmt.Println(path)

}

func ensureJSON(filepath string) {

	_, err := os.Stat(filepath)
	if err == nil {
		return
	}

	if !os.IsNotExist(err) {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	newFile, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	defer newFile.Close()

	_, err = newFile.WriteString("{}")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

}

func setDirectoryVar(data cmdArgs) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	data.allPaths[data.cmd[1]] = path

	jsonData, err := json.MarshalIndent(data.allPaths, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		os.Exit(1)
	}

	file, err := os.Create(data.jsonPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		os.Exit(1)
	}
}

func displayAllPaths(data cmdArgs) {

    fmt.Println("\n")
    for k, v := range data.allPaths {

        fmt.Printf("%v: %v\n", k, v)
    
    }
    fmt.Println("\n")

}

type cmdArgs struct {
	cmd      []string
	allPaths map[string]string
	jsonPath string
}

// map of available commands
var availCmds = map[string]func(data cmdArgs){
	"to":  changeDirectory,
	"set": setDirectoryVar,
	"ls":  displayAllPaths,
}

func main() {

	// read in json file
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	jsonDirPath := filepath.Dir(exePath)
	jsonPath := jsonDirPath + "\\fastTravel.json"
	ensureJSON(jsonPath)
	allPaths := readMap(jsonPath)

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
		jsonPath: jsonPath,
	}

	exeCmd(data)

}
