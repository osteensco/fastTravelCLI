package ft

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"ft/helpers" // Add this import
)

func PassCmd(args []string) (*Cmd, error) {
	// Dissect provided args
	cmd := ParseArgs(&args)

	// all commands are expected to lead with "-"
	// the only command that doesn't is a change directory
	// commands that are symbols have this leader added to avoid logic
	// meant for a directory path
	switch cmd.Cmd {
	case "]", "[", "..", "-":
		cmd.Cmd = fmt.Sprintf("-%s", cmd.Cmd)
	}

	// keys and key evals are given a leader of "_"
	// to help identify the appropriate function
	// in the map
	if cmd.Cmd == "" {
		cmd.Cmd = "_"
	}

	_, ok := AvailCmds[cmd.Cmd]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%v is not a valid command, use 'ft -h' or 'ft -help' for a list of valid commands", cmd.Cmd))
	}

	// We don't care about minimum number of args if we're just getting help docs for a command
	if cmd.Flags.H {
		return cmd, nil
	}

	// verify user provided correct minimum number of arguments
	// too many args will work, any args beyond expected number are simply ignored
	switch cmd.Cmd {
	case "-ls", "-]", "-[", "-..", "--", "-hist", "-help", "-h", "-version", "-v", "-is", "-update", "-u":
		break
	case "-rn", "-edit":
		if len(cmd.Args) < 2 {
			return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, see ft -help for more info", args[1:]))
		}
	default:
		if len(cmd.Args) == 0 {
			return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft <command> <path/key>", args[1:]))
		}
	}

	return cmd, nil
}

// Takes a key or relative path and returns it's absolute path or an error.
// Path is a seperate argument to handle evaluating many paths as part of a loop.
func evalPath(data *CmdAPI, path *string) (string, error) {

	// TODO
	//  - check for SPECIFIC errors in directory checks (ex: *PathError)
	//  - ensure other types of errors are caught and returned early

	var key string
	provided_string := *path

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
				return fmt.Sprintf("%s", path), nil
			}
		}

		return "", errors.New(fmt.Sprintf(InvalidDirectoryMsg, provided_string, path))

	} else {

		key = provided_string
		// handles key lookup
		p, ok := data.allPaths[key]
		if !ok {
			// handles releative directory in CWD
			dir, err := os.Stat(key)
			if err == nil {
				if dir.IsDir() {
					// in this case key is assumed to be a relative path
					p, err = filepath.Abs(key)
					if err != nil {
						return "", err
					}
					return fmt.Sprintf("%s", p), nil
				}
			}

			// handles CDPATH
			cdpath := os.Getenv("CDPATH")
			if cdpath != "" {
				cdpaths := strings.Split(cdpath, ":")
				for _, path := range cdpaths {
					cdPathResult := filepath.Join(path, key)
					dir, err := os.Stat(cdPathResult)
					if err == nil && dir.IsDir() {
						return fmt.Sprintf("%s", cdPathResult), nil
					}
				}
			}

			return "", errors.New(fmt.Sprintf(UnrecognizedKeyMsg, key))

		}

		return fmt.Sprintf("%s", p), nil

	}
}

// changeDirectory can handle key lookup, relative paths, directories in CDPATH, and key evaluation.
func changeDirectory(data *CmdAPI) error {

	path, err := evalPath(data, &data.cmd.Args[0])
	if err != nil {
		return err
	}

	// Expand and clean the path before returning it for history or actual directory change.
	// This ensures paths stored in history are always absolute and fully resolved.
	expandedPath, err := helpers.ExpandAndCleanPath(path)
	if err != nil {
		return fmt.Errorf("failed to expand and clean path '%s': %w", path, err)
	}

	// The shell wrapper reads this output to perform the actual 'cd' command.
	// By printing the expandedPath, we ensure the shell acts on a fully resolved path.
	fmt.Println(expandedPath)
	return nil // Return nil if successful, as the path has been printed.
}

func setDirectoryVar(data *CmdAPI) error {

	// reverse allPaths for easier lookup by directory instead of key
	dirs := make(map[string]string, len(data.allPaths))
	for k, v := range data.allPaths {
		dirs[v] = k
	}

	// handle setting multiple keys at once
	for _, arg := range data.cmd.Args {

		var key string
		var initialPath string // Temporary variable to hold path before expansion
		var err error

		// when key is not explicitly set to a directory, assume user is trying to set key to CWD
		// key, value pairs are expected to be in the format key=value space delimited
		if strings.Contains(arg, "=") {
			pair := strings.SplitN(arg, "=", 2) // Use SplitN to handle paths with '=' within the value
			if len(pair) != 2 {
				return errors.New("invalid key-value pair format: expected 'key=value'")
			}
			key, initialPath = pair[0], pair[1]
			initialPath, err = evalPath(data, &initialPath)
		} else {
			key = arg
			initialPath, err = os.Getwd()
		}
		if err != nil {
			return err
		}

		// Expand and clean the path before storing it, to ensure consistency in data.allPaths.
		path, err := helpers.ExpandAndCleanPath(initialPath)
		if err != nil {
			return fmt.Errorf("failed to expand and clean path '%s': %w", initialPath, err)
		}

		// verify if path is already saved to another key
		k, ok := dirs[path]
		if ok {
			var res string
			// if force setting we don't need to read in a response from the user
			if !data.cmd.Flags.Y {
				fmt.Printf(PathAlreadyExistsMsg, path, k, key)
				_, err := fmt.Fscan(data.rdr, &res)
				if err != nil {
					return err
				}
			}

			if overwrite, err := verifyInput(res, data.cmd.Flags.Y); !overwrite {
				if err != nil {
					return err
				}
				fmt.Printf(AbortedOverwriteKeyMsg, key)
				continue // Use continue to proceed to the next argument
			} else {
				if err != nil {
					return err
				}
			}
			delete(data.allPaths, k)
		}

		// verify if key is already in use
		val, ok := data.allPaths[key]
		if ok {
			var res string
			if !data.cmd.Flags.Y {
				// capture user response and act accordingly
				fmt.Printf(KeyAlreadyExistsMsg, key, val, key)
				_, err := fmt.Fscan(data.rdr, &res)
				if err != nil {
					return err
				}
			}
			if overwrite, err := verifyInput(res, data.cmd.Flags.Y); !overwrite {
				if err != nil {
					return err
				}
				fmt.Printf(AbortedOverwriteKeyMsg, key)
				continue // Use continue to proceed to the next argument
			} else {
				if err != nil {
					return err
				}
			}

		}

		// key doesn't exist yet or user wants to overwrite
		data.allPaths[key] = path
		dataUpdate(data.allPaths, data.file)
		fmt.Printf(AddKeyMsg, key, path)

	}

	return nil
}

func displayAllPaths(data *CmdAPI) error {
	fmt.Println("")
	printMap(data.allPaths, "%v: %v\n")
	fmt.Println("")
	return nil
}

func removeKey(data *CmdAPI) error {
	var res string
	key := data.cmd.Args[0]

	_, ok := data.allPaths[key]
	if !ok {
		fmt.Printf(KeyDoesNotExistMsg, key)
		return nil
	}
	if !data.cmd.Flags.Y {
		fmt.Printf(VerifyRemoveMsg, key)
		_, err := fmt.Fscan(data.rdr, &res)
		if err != nil {
			return err
		}
		if rm, err := verifyInput(res, false); !rm {
			if err != nil {
				return err
			}
			fmt.Printf(AbortRemoveKeyMsg, key)
			return nil
		} else {
			if err != nil {
				return err
			}
		}
	}
	delete(data.allPaths, key)
	dataUpdate(data.allPaths, data.file)
	fmt.Printf(RemoveKeyMsg, key)
	return nil
}

func renameKey(data *CmdAPI) error {
	originalKey := data.cmd.Args[0]
	newKey := data.cmd.Args[1]

	// Check if new key already exists
	_, ok := data.allPaths[newKey]
	if ok {
		return fmt.Errorf(RenameKeyExistsMsg, newKey)
	}

	// Check if original key exists
	val, ok := data.allPaths[originalKey]
	if !ok {
		return fmt.Errorf(KeyDoesNotExistMsg, originalKey)
	}

	// Perform rename
	delete(data.allPaths, originalKey)
	data.allPaths[newKey] = val
	dataUpdate(data.allPaths, data.file)
	fmt.Printf(RenameKeyMsg, originalKey, newKey)
	return nil
}