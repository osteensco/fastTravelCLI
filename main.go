package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)


// TODO
// - tests
//      - main() test should loop over every command, not just help
//      - install test should automagically run cleanup.sh
// - rm cmd
//      - this cmd should prompt the user to confirm before removing
//          - probably for rn cmd as well
//      - find way to add user interactivity between child process and parent process
//          - try using 'trap' in bash to trigger commands based on exe output
// - set cmd
//      - should tell you if you already have the dir saved
//      - should tell you if you are about to overwrite a key
// - rn cmd
//      - should display error when attempting to rename key that doesn't exist
// - features:
//      - ft swap [key1] [key2]
//          - swaps the dirs of the two keys given
//      - ft ?
//          - ask fastTravel if the curr dir is saved
//      - remember prev n directories? n=10?15?20?
//      - add a spinner?



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



func readMap(file *os.File) map[string]string {

    pathMap := make(map[string]string)
    

    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Println("Error getting file info: ", err)
        os.Exit(1)
    }

    // take file size in bytes and make a buffer of that size
    size := fileInfo.Size()
    buff := make([]byte, size)

    // read entire file into memory
    _, err = file.Read(buff)
    if err != nil {
        fmt.Println("Error reading file into buffer: ", err)
        os.Exit(1)
    }

    // key length integer should always fit in 1 byte
    var keyLen uint8
    // value length integer should always fit in 2 bytes
    var valLen uint16
    // sliding pointer to navigate buffer
    var offset uint
    
    // iterate through buffer and deserialize
    for offset < uint(len(buff)) {

        // read length of key, use length to read in key, adjust offset
        // simple type conversion since length is only 1 byte and not a []byte
        keyLen = uint8(buff[offset])
        offset++
        kl := uint(keyLen)
        keyBytes := buff[offset:offset+kl]
        offset += kl
        
        // read length of value, use length to read in value, adjust offset
        // length contained in 2 bytes, nedd to convert []byte to a uint16 value
        valLen = binary.LittleEndian.Uint16(buff[offset:offset+2])
        offset += 2
        vl := uint(valLen)
        valBytes := buff[offset:offset+vl]
        offset += vl
        // add key-value to map
        pathMap[string(keyBytes)] = string(valBytes)

    }

	return pathMap

}



func ensureData(filepath string) *os.File {

    file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
    if err != nil {
        fmt.Println("Error opening file: ", err)
        os.Exit(1)
    }

    return file

}




func dataUpdate(hashmap map[string]string, file *os.File) {
    
    var buffer []byte
    for key, val := range hashmap {

        keyBytes := []byte(key)
        valBytes := []byte(val)
        
        keyLen := make([]byte, 1)
        keyLen[0] = byte(uint8(len(keyBytes)))
        
        valLen := make([]byte,2)
        binary.LittleEndian.PutUint16(valLen, uint16(len(valBytes)))

        // create an array of array of bytes for optimal concatenation
        allBytes := [][]byte{keyLen, keyBytes, valLen, valBytes}

        // get the length for key-value pair to allocate memory
        var pairLen int
        for _, s := range allBytes {
            pairLen += len(s)
        }
        // create new slice and append all []byte from allBytes
        pair := make([]byte, pairLen)
        var i int
        for _, s := range allBytes {
            i += copy(pair[i:], s)
        }
        // append completed pair to buffer []byte
        allPairs := make([]byte, len(buffer)+len(pair))
        copy(allPairs, buffer)
        copy(allPairs[len(buffer):], pair)
        buffer = allPairs
    }
    
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

    _, err = file.Write(buffer)
    if err != nil {
        fmt.Println("Error writing contents of buffer to file: ", err)
        os.Exit(1)
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

    dataUpdate(data.allPaths, data.file)
    fmt.Printf("Added destination %v", data.cmd[1] )

}

func displayAllPaths(data cmdArgs) {

    printMap(data.allPaths)

}

func removeKey(data cmdArgs) {
    
    delete(data.allPaths, data.cmd[1])
    dataUpdate(data.allPaths, data.file)
    fmt.Printf("Removed '%v' destination", data.cmd[1])
    
}

func renameKey(data cmdArgs) {
    
    originalKey := data.cmd[1]
    newKey := data.cmd[2]
    path, ok := data.allPaths[originalKey]
    if !ok {
        err := errors.New(fmt.Sprintf("Cannot rename %v, key does not exist. Run 'ft ls' to see all keys.", originalKey))
        fmt.Println(err)
        return
   }
    delete(data.allPaths, originalKey)
    data.allPaths[newKey] = path 

    dataUpdate(data.allPaths, data.file)
    fmt.Printf("%v renamed to %v", originalKey, newKey)

}

func showHelp(data cmdArgs) {
    
    printMap(cmdDesc)

}

type cmdArgs struct {
	cmd      []string
	allPaths map[string]string
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
	}

	exeCmd(data)

}
