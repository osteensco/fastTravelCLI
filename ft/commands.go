package ft

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func PassCmd(args []string) ([]string, error) {
	cmd := args[1]

	// all commands are expected to lead with "-"
	// the only command that doesn't is a change directory
	// commands that are symbols have this leader added to avoid logic
	// meant for a directory path
	switch cmd {
	case "]", "[", "..", "-":
		cmd = fmt.Sprintf("-%s", cmd)
	default:
		break
	}

	// keys and key evals are given a leader of "_"
	// to help identify the appropriate function
	// in the map
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
		return nil, errors.New(fmt.Sprintf("%v is not a valid command, use 'ft -h' or 'ft -help' for a list of valid commands", cmd))
	}

	// verify user provided correct minimum number of arguments
	// too many args will work, any args beyond expected number are simply ignored
	switch cmd {
	case "-ls", "-]", "-[", "-..", "--":
		return []string{cmd}, nil
	// providing help for a specific command may be needed in the future
	case "-help", "-h", "-version", "-v":
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

func changeDirectory(data *CmdArgs) error {
	if len(data.allPaths) == 0 {
		fmt.Printf("No fast travel locations set, set locations by navigating to desired destination directory and using 'ft -set <key>' \n")
		return nil
	}

	var key string
	provided_string := data.cmd[1]

	if strings.Contains(provided_string, "/") {

		path_array := strings.Split(provided_string, "/")
		eval_array := make([]string, len(path_array))

		// key evaluation if the first string before "/" delimeter is key
		key = path_array[0]
		if p, ok := data.allPaths[key]; ok {
			eval_array[0] = p
		} else {
			eval_array[0] = key
		}
		for i, str := range path_array {
			if i != 0 {
				eval_array[i] = str
			}
		}

		// handles evaluated path and relative paths
		path := strings.Join(eval_array, "/")
		dir, err := os.Stat(path)
		if err == nil {
			if dir.IsDir() {
				fmt.Println(path)
				return nil
			}
		}

		fmt.Printf("Provided path '%s' evaluates to '%s' which is not a valid directory. Use 'ft -ls' to see all saved destinations. \n", provided_string, path)
		return nil

	} else {

		key = provided_string
		p, ok := data.allPaths[key]
		if !ok {
			fmt.Printf("Did not recognize key '%s', use 'ft -ls' to see all saved destinations. If this is a relative path use './%s' or '%s/'. \n", key, key, key)
			return nil
		}

		fmt.Println(p)
		return nil

	}
}

func setDirectoryVar(data *CmdArgs) error {
	key := data.cmd[1]
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	for k, v := range data.allPaths {
		if path == v {
			fmt.Printf("Path '%s' already exists with key '%s', overwrite key '%s' \n", path, k, k)
			var res string
			_, err := fmt.Fscan(data.rdr, &res)
			if err != nil {
				return err
			}
			if overwrite, err := verifyInput(res); !overwrite {
				if err != nil {
					return err
				}
				fmt.Printf("Aborted overwriting of key '%s' \n", key)
				return nil
			} else {
				if err != nil {
					return err
				}
			}
			delete(data.allPaths, k)
			data.allPaths[key] = path
			dataUpdate(data.allPaths, data.file)
			fmt.Printf("Renamed key '%s' to '%s' whose value is '%s' \n", k, key, path)
			return nil
		}
	}

	val, ok := data.allPaths[key]
	if ok {
		// capture user response and act accordingly
		fmt.Printf("Key '%s' already exists with value '%s', overwrite key '%s'?(y/n) ", key, val, key)
		var res string
		_, err := fmt.Fscan(data.rdr, &res)
		if err != nil {
			return err
		}
		if overwrite, err := verifyInput(res); !overwrite {
			if err != nil {
				return err
			}
			fmt.Printf("Aborted overwriting of key '%s' \n", key)
			return nil
		} else {
			if err != nil {
				return err
			}
		}

	}

	// key doesn't exist yet or user wants to overwrite
	data.allPaths[key] = path
	dataUpdate(data.allPaths, data.file)
	fmt.Printf("Added destination '%s': '%s' \n", key, path)
	return nil
}

func displayAllPaths(data *CmdArgs) error {
	printMap(data.allPaths)
	return nil
}

func removeKey(data *CmdArgs) error {
	var res string
	key := data.cmd[1]

	_, ok := data.allPaths[key]
	if !ok {
		fmt.Printf("Key '%s' does not exist. Run 'ft -ls' to see all keys. \n", key)
		return nil
	}
	fmt.Printf("Are you sure you want to remove the key '%s'? (y/n) ", key)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		return err
	}
	if rm, err := verifyInput(res); !rm {
		if err != nil {
			return err
		}
		fmt.Printf("Aborted removal of key '%s' \n", key)
		return nil
	} else {
		if err != nil {
			return err
		}
	}
	delete(data.allPaths, key)
	dataUpdate(data.allPaths, data.file)
	fmt.Printf("Removed '%s' destination \n", key)
	return nil
}

func renameKey(data *CmdArgs) error {
	originalKey := data.cmd[1]
	newKey := data.cmd[2]

	_, ok := data.allPaths[newKey]
	if ok {
		fmt.Printf("Key '%s' already exists, please choose something else. \n", newKey)
		return nil
	}
	path, ok := data.allPaths[originalKey]
	if !ok {
		fmt.Printf("Cannot rename '%s', key does not exist. Run 'ft -ls' to see all keys. \n", originalKey)
		return nil
	}

	var res string

	fmt.Printf("Are you sure you want to rename the key '%s' to '%s'? (y/n) ", originalKey, newKey)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		return err
	}
	if rm, err := verifyInput(res); !rm {
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf("Aborted renaming of key '%s' to '%s'. \n", originalKey, newKey)
		return nil
	} else {
		if err != nil {
			return err
		}
	}
	delete(data.allPaths, originalKey)
	data.allPaths[newKey] = path

	dataUpdate(data.allPaths, data.file)
	fmt.Printf("Key '%s' renamed to '%s'. \n", originalKey, newKey)
	return nil
}

func showHelp(data *CmdArgs) error {
	printMap(CmdDesc)
	return nil
}

func showVersion(data *CmdArgs) error {
	fmt.Println(`
     __           _  _____                     _   ___   __   _____ - -  -  -   -   -
    / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ - -  -  -   -   -
   | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/  - -  -   -   -
   |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  - -  -  -   -   -
   |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/ - -  -  -   -   -

		`)
	fmt.Println("version:\t", Version)
	return nil
}

// Used for commands that are simply handled by the shell function
func passToShell(data *CmdArgs) error {
	c := data.cmd[0]
	command := string(c[1:])

	switch command {
	case "]", "[", "..", "-":
		fmt.Println(command)
	default:
		return errors.New(fmt.Sprintf("Tried to pass command to shell, but '%s' is not a valid command for the shell function.", command))
	}

	return nil
}
