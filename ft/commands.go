package ft

import (
	"errors"
	"fmt"
	"os"
)

func PassCmd(args []string) ([]string, error) {

	cmd := args[1]

	// all commands are expected to lead with "-"
	// the only command that doesn't is a CD to provided key
	if string(cmd[0]) != "-" {
		cmd = "_"
		parsedCmd := make([]string, 3)
		parsedCmd[0] = args[0]
		parsedCmd[1] = cmd
		parsedCmd[2] = args[1]
		args = parsedCmd
	}
	_, ok := AvailCmds[cmd]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v is not a valid command, use 'ft help' for valid commands", cmd))
	}

	// verify user provided correct minimum number of arguments
	// too many args will work, any args beyond expected number are simply ignored
	switch cmd {
	case "-ls":
		break
	// providing help for a specific command may be needed in the future
	case "-help", "-h":
		break
	case "-rn":
		if len(args) <= 3 {
			return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft rn <key> <newKey>", args[1:]))
		}
	default:
		if len(args) <= 2 {
			return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft <command> <path/key>", args[1:]))
		}
	}

	// return args without 'ft'
	return args[1:], nil

}

func changeDirectory(data *CmdArgs) {

	if len(data.allPaths) == 0 {
		fmt.Printf("No fast travel locations set, set locations by navigating to desired destination directory and using 'ft -set <key>' ")
		os.Exit(1)
	}

	key := data.cmd[1]

	p, ok := data.allPaths[key]
	if !ok {
		fmt.Printf("Attempt to fast travel to key %s failed! Use 'ft -ls' to see all saved destinations.", key)
		os.Exit(1)
	}

	fmt.Println(p)

}

func setDirectoryVar(data *CmdArgs) {

	key := data.cmd[1]
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	for k, v := range data.allPaths {
		if path == v {
			fmt.Printf("Path '%s' already exists with key %s", path, k)
			return
		}
	}

	val, ok := data.allPaths[key]
	if ok {
		// capture user response and act accordingly
		fmt.Printf("Key '%s' already exists with value %s, overwrite key '%s'?(y/n)", key, val, key)
		var res string
		_, err := fmt.Fscan(data.rdr, &res)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if overwrite, err := verifyInput(res); !overwrite {
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			fmt.Printf("Aborted overwriting of key %v", key)
			os.Exit(1)
		} else {
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
		}

	}

	// key doesn't exist yet or user wants to overwrite
	data.allPaths[key] = path
	dataUpdate(data.allPaths, data.file)
	fmt.Printf("Added destination %s", data.cmd[1])

}

func displayAllPaths(data *CmdArgs) {

	printMap(data.allPaths)

}

func removeKey(data *CmdArgs) {

	var res string
	key := data.cmd[1]

	fmt.Printf("Are you sure you want to remove the key '%v'? (y/n)", key)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if rm, err := verifyInput(res); !rm {
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Printf("Aborted removal of key %v", key)
		os.Exit(1)
	} else {
		if err != nil {
			fmt.Println("Error: ", err)
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

	fmt.Printf("Are you sure you want to rename the key '%v'? (y/n)", newKey)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if rm, err := verifyInput(res); !rm {
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("Aborted renaming of key %v", newKey)
		os.Exit(1)
	} else {
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}
	delete(data.allPaths, originalKey)
	data.allPaths[newKey] = path

	dataUpdate(data.allPaths, data.file)
	fmt.Printf("%v renamed to %v", originalKey, newKey)

}

func showHelp(data *CmdArgs) {

	printMap(CmdDesc)

}
