package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
    "strings"
)

// TODO
// - rm cmd
     // - cannot do check prior to deletion since we are passing STDout to a bash function (callstack, duh)
     // - need to delete in json as well, not just in memory
         // - need to create a json update function at this point
// - check if valid command, display valid command response, display insuff args response if valid cmd with insuff args
 // - investigate data persistence alternatives



func printMap(hashmap map[string]string) {

    fmt.Println("")

    keys := make([]string, 0, len(hashmap))
    for k := range hashmap {
        keys = append(keys, k)
    }

    sort.Strings(keys)

    for i := range keys {

        fmt.Printf("%v: %v\n", keys[i], hashmap[keys[i]])
    
    }

    fmt.Println("")

}

func passCmd(args []string) ([]string, error) {
    
    switch args[1] {
        case "ls":
            break
        case "help":
            break
        default:
            if len(args) <= 2 {
                return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft <command> <path/key>", args[1:]))
            }
    }

    // if cmd := args[1]; cmd != "ls" || cmd != "help" {
    //     if len(args) <= 2 {
    //         return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft <command> <path/key>", args[1:]))
    //     }
    // }
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
		fmt.Printf("No fast travel locations set, set locations by navigating to desired destination directory and using 'ft set <key>' ")
		os.Exit(1)
	}

	path := data.allPaths[data.cmd[1]]
	distro := os.Getenv("WSL_DISTRO_NAME")
	if len(distro) == 0 {

    	prefix := "/mnt/"
        drive := strings.Split(path, ":")[0]
        path = strings.Replace(path, drive, strings.ToLower(drive), 1)
        path = strings.Replace(path, ":", "", 1)
    	path = strings.Replace(path, "\\", "/", -1)
    	path = prefix + path
        // path = strings.ToLower(prefix + path)

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

    printMap(data.allPaths)

}

func removeKey(data cmdArgs) {
    
    reader := bufio.NewReader(os.Stdin)

        fmt.Printf("Are you sure you want to delete '%v: %v'? y/n\n", data.cmd[1], data.allPaths[data.cmd[1]])
        resp, _ := reader.ReadString('\n')
        resp = strings.Trim(resp, "\n")

        switch strings.ToLower(resp) {
            case "y": 
                delete(data.allPaths, data.cmd[1])
                fmt.Printf("Removed '%v' destination", data.cmd[1])
                return
            case "n":
                fmt.Printf("Did not remove '%v: %v'", data.cmd[1], data.allPaths[data.cmd[1]])
                return
            default:
                fmt.Printf("'%v' is not a valid response, type y or n", resp)
        }   


    
}

func renameKey(data cmdArgs) {
    
    originalKey := data.cmd[1]
    newKey := data.cmd[2]
    path := data.allPaths[originalKey]
    delete(data.allPaths, originalKey)
    data.allPaths[newKey] = path 

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

    fmt.Printf("%v renamed to %v", originalKey, newKey)

}

func showHelp(data cmdArgs) {
    
    printMap(cmdDesc)

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
