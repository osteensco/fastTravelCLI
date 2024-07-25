package main

import (
	// "bufio"
	"encoding/binary"
	// "encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	// "time"
)

// TODO
// - rm cmd
//      - this cmd should prompt the user to confirm before removing
//          - probably for rn cmd as well
//      - find way to add user interactivity between child process and parent process
//          - try using 'trap' in bash to trigger commands based on exe output
// - investigate data persistence alternatives
//      - need to test json vs bin files, at maybe 10mb?
// - explore more efficient ways of reading in data
//      - should entire file be read into memory one every use?
//      - what would be the alternative? don't want to use a db, trying to avoid dependencies.
// - rename json functions to 'data'




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
    
    cmd := args[1]

    _, ok := availCmds[cmd]
    if !ok {
        return nil, errors.New(fmt.Sprintf("%v is not a valid command, use 'ft help' for valid commands", cmd))
    }
    switch cmd {
        case "ls":
            break
        case "help":
            break
        case "rn":
            if len(args) <= 3 {
                return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft rn <key> <newKey>", args[1:]))
            }
        default:
            if len(args) <= 2 {
                return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft <command> <path/key>", args[1:]))
            }
    }

	return args[1:], nil

}

// func readMap(jsonPath string) map[string]string {
// 	file, err := os.Open(jsonPath)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		os.Exit(1)
// 	}

// 	defer file.Close()

// 	var pathMap map[string]string
// 	decoder := json.NewDecoder(file)
// 	if err := decoder.Decode(&pathMap); err != nil {
// 		fmt.Println("Error:", err)
// 		os.Exit(1)
// 	}

// 	return pathMap

// }


func readMap(file *os.File) map[string]string {

    pathMap := make(map[string]string)

    _, err := file.Seek(0,0) 
    if err != nil {
        fmt.Println("Error seeking to beginning of file: ", err)
        os.Exit(1)
    }

    // key length integer should always fit in 8 bits
    var keyLen uint8
    // value length integer should always fit in 16 bits
    var valLen uint16

    // this should iterate until the end of the file
    for {
        // read length of key, use length to read in key
        err = binary.Read(file, binary.LittleEndian, &keyLen)
        // End Of File error should only ever happen here
        if err == io.EOF {
            break
        }
        if err != nil {
            fmt.Println("Error reading keyLen binary: ", err)
            os.Exit(1)
        }
        keyBytes := make([]byte, keyLen)
        _, err = file.Read(keyBytes)
        if err != nil {
            fmt.Println("Error reading in keyBytes: ", err)
        }


        // read length of value, use length to read in value
        err = binary.Read(file, binary.LittleEndian, &valLen)
        if err != nil {
            fmt.Println("Error reading valLen binary: ", err)
            os.Exit(1)
        }
        valBytes := make([]byte, valLen)
        _, err = file.Read(valBytes)
        if err != nil {
            fmt.Println("Error reading in valBytes: ", err)
        }


        pathMap[string(keyBytes)] = string(valBytes)

    }

	return pathMap

}


// func ensureJSON(filepath string) {

// 	_, err := os.Stat(filepath)
// 	if err == nil {
// 		return
// 	}

// 	if !os.IsNotExist(err) {
// 		fmt.Println("Error: ", err)
// 		os.Exit(1)
// 	}

// 	newFile, err := os.Create(filepath)
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		os.Exit(1)
// 	}

// 	defer newFile.Close()

// 	_, err = newFile.WriteString("{}")
// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 		os.Exit(1)
// 	}

// }


func ensureJSON(filepath string) *os.File {

    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        fmt.Println("Error opening file: ", err)
        os.Exit(1)
    }

    return file

}


// func jsonUpdate(hashmap map[string]string, filepath string) {
    
// 	jsonData, err := json.MarshalIndent(hashmap, "", "  ")
// 	if err != nil {
// 		fmt.Println("Error marshalling JSON:", err)
// 		os.Exit(1)
// 	}

// 	file, err := os.Create(filepath)
// 	if err != nil {
// 		fmt.Println("Error creating file:", err)
// 		os.Exit(1)
// 	}

// 	defer file.Close()

// 	_, err = file.Write(jsonData)
// 	if err != nil {
// 		fmt.Println("Error writing JSON to file:", err)
// 		os.Exit(1)
// 	}

// }



func jsonUpdate(hashmap map[string]string, file *os.File) {
    
    err := file.Truncate(0) 
    if err != nil {
        fmt.Println("Error truncating file: ", err)
        os.Exit(1)
    }
    
    _, err = file.Seek(0,0) 
    if err != nil {
        fmt.Println("Error seeking to beginning of file: ", err)
        os.Exit(1)
    }


    for key, val := range hashmap {

        keyBytes := []byte(key)
        valBytes := []byte(val)
        keyLen := uint8(len(keyBytes))
        valLen := uint16(len(valBytes))

        
        err = binary.Write(file, binary.LittleEndian, keyLen)
        if err != nil {
            fmt.Println("Error writing keyLen binary: ", err)
            os.Exit(1)
        }
        _, err = file.Write(keyBytes)
        if err != nil {
            fmt.Println("Error writing keyBytes binary: ", err)
            os.Exit(1)
        }

        err = binary.Write(file, binary.LittleEndian, valLen)
        if err != nil {
            fmt.Println("Error writing valLen binary: ", err)
            os.Exit(1)
        }
        _, err = file.Write(valBytes)
        if err != nil {
            fmt.Println("Error writing valBytes binary: ", err)
            os.Exit(1)
        }

    }

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

	}

    fmt.Println(path)

}

func setDirectoryVar(data cmdArgs) {

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	data.allPaths[data.cmd[1]] = path

    // jsonUpdate(data.allPaths, data.jsonPath)
    jsonUpdate(data.allPaths, data.file)
    fmt.Printf("Added destination %v", data.cmd[1] )

}

func displayAllPaths(data cmdArgs) {

    printMap(data.allPaths)

}

func removeKey(data cmdArgs) {
    
    delete(data.allPaths, data.cmd[1])
    // jsonUpdate(data.allPaths, data.jsonPath)
    jsonUpdate(data.allPaths, data.file)
    fmt.Printf("Removed '%v' destination", data.cmd[1])
    
}

func renameKey(data cmdArgs) {
    
    originalKey := data.cmd[1]
    newKey := data.cmd[2]
    path := data.allPaths[originalKey]
    delete(data.allPaths, originalKey)
    data.allPaths[newKey] = path 

    // jsonUpdate(data.allPaths, data.jsonPath)
    jsonUpdate(data.allPaths, data.file)
    fmt.Printf("%v renamed to %v", originalKey, newKey)

}

func showHelp(data cmdArgs) {
    
    printMap(cmdDesc)

}

type cmdArgs struct {
	cmd      []string
	allPaths map[string]string
	// jsonPath string
    file *os.File
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

// func measureSpeed(start time.Time) {
//     fmt.Printf("execution took %v", time.Since(start))
// }

func main() {

    // defer measureSpeed(time.Now())
	
    // read in json file
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	jsonDirPath := filepath.Dir(exePath)
	// jsonPath := jsonDirPath + "\\fastTravel.json"
	jsonPath := jsonDirPath + "\\fastTravel.bin"
    file := ensureJSON(jsonPath)
    defer file.Close()
	// allPaths := readMap(jsonPath)
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
		// jsonPath: jsonPath,
        file: file,
	}

	exeCmd(data)

}
