package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)


func passCmd(args []string) ([]string, error) {

    if len(args) <= 2 {
        return nil, errors.New("Insufficient args provided, usage: ftrav <command> <path/key>")
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


func changeDirectory(cmd []string, allPaths map[string]string, jsonPath string) {

    path := allPaths[cmd[1]]

    err := os.Chdir(path)
    if err != nil {
        fmt.Printf("Fast travel failed! %v", err)
        os.Exit(1)
    }

    fmt.Println("moved to ", path)

}

func setDirectoryVar(cmd []string, allPaths map[string]string, jsonPath string) {
     
    path, err := os.Getwd()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    allPaths[cmd[2]] = path
    
    jsonData, err := json.MarshalIndent(allPaths, "", "  ")
    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        os.Exit(1)
    }
 
    file, err := os.Create(jsonPath)
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

// func displayAllPaths(cmd []string, allPaths map[string]string, jsonPath string) {}
// type cmdArgs struct {
//     cmd []string 
//     allPaths map[string]string  
//     jsonPath string 
// }

func main() {
    
    // read in json file
    jsonPath := "fastTravel.json"
    allPaths := readMap(jsonPath)

    // map available ftrav commands 
    availCmds := map[string]func(cmd []string, allPaths map[string]string, jsonPath string) {
        "to": changeDirectory,
        "set": setDirectoryVar,
        // "ls": displayAllPaths
    }
    
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

    exeCmd(inputCommand, allPaths, jsonPath)

}



