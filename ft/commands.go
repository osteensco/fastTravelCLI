package ft

import (
	"errors"
	"fmt"
	"os"
)

func PassCmd(args []string) ([]string, error) {

	cmd := args[1]

	_, ok := AvailCmds[cmd]
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

func changeDirectory(data *CmdArgs) {

	if len(data.allPaths) == 0 {
		fmt.Printf("No fast travel locations set, set locations by navigating to desired destination directory and using 'ft set <key>' ")
		os.Exit(1)
	}

	p := data.allPaths[data.cmd[1]]
	path := SanitizeDir(p)
	fmt.Println(path)

}

func setDirectoryVar(data *CmdArgs) {

	key := data.cmd[1]
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for _, v := range data.allPaths {
		if path == v {
			fmt.Printf("Path '%v' already exists", path)
			return
		}
	}

	_, ok := data.allPaths[key]
	if !ok {
		data.allPaths[key] = path
		dataUpdate(data.allPaths, data.file)
		fmt.Printf("Added destination %v", data.cmd[1])
	} else {
		fmt.Printf("Key '%v' already exists", key)
	}

}

func displayAllPaths(data *CmdArgs) {

	printMap(data.allPaths)

}

func removeKey(data *CmdArgs) {

	var res string
	key := data.cmd[1]

	fmt.Printf("Are you sure you want to remove the key '%v'?", key)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if rm, err := verifyInput(res); !rm {
		fmt.Printf("Aborted removal of key %v", key)
		os.Exit(1)
	} else {
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
	delete(data.allPaths, key)
	dataUpdate(data.allPaths, data.file)
	fmt.Printf("Removed '%v' destination", key)

}

func renameKey(data *CmdArgs) {

	originalKey := data.cmd[1]
	newKey := data.cmd[2]

	_, ok := data.allPaths[newKey]
	if ok {
		fmt.Printf("Key %v already exists, please choose something else.", newKey)
		os.Exit(1)
	}
	path, ok := data.allPaths[originalKey]
	if !ok {
		err := errors.New(fmt.Sprintf("Cannot rename %v, key does not exist. Run 'ft ls' to see all keys.", originalKey))
		fmt.Println(err)
		os.Exit(1)
	}

	var res string

	fmt.Printf("Are you sure you want to rename the key '%v'?", newKey)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if rm, err := verifyInput(res); !rm {
		fmt.Printf("Aborted renaming of key %v", newKey)
		os.Exit(1)
	} else {
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
	delete(data.allPaths, originalKey)
	data.allPaths[newKey] = path

	dataUpdate(data.allPaths, data.file)
	fmt.Printf("%v renamed to %v", originalKey, newKey)

}

func showHelp(data *CmdArgs) {

	printMap(cmdDesc)

}
