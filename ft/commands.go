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
	case "-ls", "-]", "-[", "-..", "--", "-hist":
		return []string{cmd}, nil
	// providing help for a specific command may be needed in the future
	case "-help", "-h", "-version", "-v", "-is", "-update", "-u":
		break
	case "-rn", "-edit":
		if len(args) <= 3 {
			return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, see ft -help for more info", args[1:]))
		}
	default:
		if len(args) <= 2 {
			return nil, errors.New(fmt.Sprintf("Insufficient args provided %v, usage: ft <command> <path/key>", args[1:]))
		}
	}

	// return args without 'ft'
	return args[1:], nil
}

// Takes a key or relative path and returns it's full explict path or an error.
func evalPath(data *CmdArgs) (string, error) {

	// TODO
	//  - check for SPECIFIC errors in directory checks (ex: *PathError)
	//  - ensure other types of errors are caught and returned early

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
				return path, nil
			}
		}

		return fmt.Sprintf(InvalidDirectoryMsg, provided_string, path), nil

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
					return p, nil
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
						return cdPathResult, nil
					}
				}
			}

			return fmt.Sprintf(UnrecognizedKeyMsg, key), nil

		}

		return p, nil

	}
}

// changeDirectory can handle key lookup, relative paths, directories in CDPATH, and key evaluation.
func changeDirectory(data *CmdArgs) error {

	path, err := evalPath(data)
	fmt.Println(path)
	return err

}

func setDirectoryVar(data *CmdArgs) error {

	// reverse allPaths for easier lookup by directory instead of key
	dirs := make(map[string]string, len(data.allPaths))
	for k, v := range data.allPaths {
		dirs[v] = k
	}

	force := false

	// handle setting multiple keys at once
	for i, arg := range data.cmd {
		if i == 0 {

			if arg[len(arg)-1] == 'f' {
				force = true
			}
			continue
		}

		var key string
		var path string
		var err error

		// if not key not explicitly set to a directory, assume user is trying to set key to CWD
		if strings.Contains(arg, "=") {
			pair := make([]string, 2)
			pair = strings.Split(arg, "=")
			key, path = pair[0], pair[1]
		} else {
			key = arg
			path, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		// verify if path is already saved to another key
		k, ok := dirs[path]
		if ok {
			var res string
			// if force setting we don't need to read in a response from the user
			if !force {
				fmt.Printf(PathAlreadyExistsMsg, path, k, key)
				_, err := fmt.Fscan(data.rdr, &res)
				if err != nil {
					return err
				}
			}

			if overwrite, err := verifyInput(res, force); !overwrite {
				if err != nil {
					return err
				}
				fmt.Printf(AbortedOverwriteKeyMsg, key)
				return nil
			} else {
				if err != nil {
					return err
				}
			}
			delete(data.allPaths, k)
			data.allPaths[key] = path
			dataUpdate(data.allPaths, data.file)
			fmt.Printf(PathOverwriteMsg, key, path)
			return nil
		}

		// verify if key is already in use
		val, ok := data.allPaths[key]
		if ok {
			var res string
			if !force {
				// capture user response and act accordingly
				fmt.Printf(KeyAlreadyExistsMsg, key, val, key)
				_, err := fmt.Fscan(data.rdr, &res)
				if err != nil {
					return err
				}
			}
			if overwrite, err := verifyInput(res, force); !overwrite {
				if err != nil {
					return err
				}
				fmt.Printf(AbortedOverwriteKeyMsg, key)
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
		fmt.Printf(AddKeyMsg, key, path)

	}

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
		fmt.Printf(KeyDoesNotExistMsg, key)
		return nil
	}
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
	delete(data.allPaths, key)
	dataUpdate(data.allPaths, data.file)
	fmt.Printf(RemoveKeyMsg, key)
	return nil
}

func renameKey(data *CmdArgs) error {
	originalKey := data.cmd[1]
	newKey := data.cmd[2]

	_, ok := data.allPaths[newKey]
	if ok {
		fmt.Printf(RenameKeyAlreadyExistsMsg, newKey)
		return nil
	}
	path, ok := data.allPaths[originalKey]
	if !ok {
		fmt.Printf(RenameKeyDoesNotExistMsg, originalKey)
		return nil
	}

	var res string

	fmt.Printf(VerifyRenameMsg, originalKey, newKey)
	_, err := fmt.Fscan(data.rdr, &res)
	if err != nil {
		return err
	}
	if rm, err := verifyInput(res, false); !rm {
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Printf(AbortRenameKeyMsg, originalKey, newKey)
		return nil
	} else {
		if err != nil {
			return err
		}
	}
	delete(data.allPaths, originalKey)
	data.allPaths[newKey] = path

	dataUpdate(data.allPaths, data.file)
	fmt.Printf(RenamedKeyMsg, originalKey, newKey, path)
	return nil
}

// Edits a saved path matching a specific prefix to the given new directory name.
// Path provided can be explicit, another key, a sub directory of a key, or relative path.
//
// Path is evaluated to explicit path to ensure replacement accuracy.
// Simple substring replacement would lead to erroneous results
// i.e. ft -edit something different
// This would update mydir/something/project1 and mydir/something/project2
// but erroneously updates myotherdir/important_project/something
// ft -edit mydir/something different
// Using explicit path (or evaluated path) avoids this problem.
func editPath(data *CmdArgs) error {

	// evaluate path to handle relative path, CDPATH, other keys, etc
	path, err := evalPath(data)
	if err != nil {
		return err
	}

	// replace directory name with new name
	newDirName := data.cmd[2]
	pathArray := strings.Split(path, "/")
	pathArray[len(pathArray)-1] = newDirName
	newPath := strings.Join(pathArray, "/")

	// update all keys that contain this prefix
	for k, v := range data.allPaths {
		if strings.Contains(v, path) {
			pathReplacement := strings.Replace(v, path, newPath, 1)
			dir, err := os.Stat(pathReplacement)
			if err != nil {
				return err
			}
			if dir.IsDir() {
				data.allPaths[k] = newPath
			} else {
				return errors.New(fmt.Sprintf("%v is not a Directory", dir))
			}
		}
	}

	return nil
}

func showHelp(data *CmdArgs) error {
	printMap(CmdDesc)
	return nil
}

func showVersion(data *CmdArgs) error {
	fmt.Print(Logo)
	fmt.Println("version:\t", Version)
	return nil
}

func showDirectoryVar(data *CmdArgs) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	pathMaps := data.allPaths

	for k := range pathMaps {
		v, _ := pathMaps[k]

		if strings.Compare(v, dir) == 0 {
			fmt.Printf(IsKeyMsg, dir, k)
			return nil
		}
	}

	fmt.Printf(IsNotKeyMsg, dir)
	return nil
}

func updateFT(data *CmdArgs) error {

	// determine if version was provided
	// default to latest if none provided
	version := "latest"
	if len(data.cmd) > 1 {
		version = data.cmd[1]
	}

	// verify/obtain version
	var endpoint string
	if version == "latest" {
		endpoint = EndpointLatestGH
	} else {
		endpoint = fmt.Sprintf(EndpointGH, version)
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending Get request to %q: %v", endpoint, err))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(
			fmt.Sprintf(
				"Error while attempting to retrieve version from github repo - status: %v %s. \n %s",
				resp.StatusCode,
				http.StatusText(resp.StatusCode),
				endpoint,
			),
		)
	}

	type respBody struct {
		TagName string `json:"tag_name"`
	}
	var body respBody
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&body)
	if err != nil {
		return err
	}
	version = body.TagName

	// verify current version is not the version attempting to be updated to
	if Version == version {
		fmt.Println("fastTravelCLI version is already ", version)
		return nil
	} else {
		fmt.Printf("fastTravelCLI %v updating to %v \n", Version, version)
	}

	// make temp directory, clone the repo
	tmpdir, err := os.MkdirTemp("", "ft_update_temp_folder")
	if err != nil {
		return errors.New(fmt.Sprintln("Error making temp directory: ", err))
	}
	defer os.RemoveAll(tmpdir)

	err = os.Chdir(tmpdir)
	if err != nil {
		return err
	}

	// TODO
	// may need to handle git clone from https, ssh, or cli

	// just using https for now
	GitCloneCMD[3] = version

	clonecmd := exec.Command(GitCloneCMD[0], GitCloneCMD[1:]...)
	err = clonecmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("Error cloning repo: %v - Command used: %q", err, clonecmd.String()))
	}
	// TODO
	// handle distinction between stable and nightly

	// run install script
	output := ""
	script := ""
	err = os.Chdir(GitCloneDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error! Could not change to dir %q", GitCloneDir))
	}

	// skip install script if function call during testing
	if UPDATEMOCK {
		return nil
	}

	switch opsys := runtime.GOOS; opsys {
	case "linux":
		script = "linux.sh"
	case "darwin":
		script = "mac.sh"
	default:
		return errors.New(fmt.Sprintf("OS %s is not handled in the update command!", opsys))
	}

	fmt.Printf("Running %s script from install folder... \n", script)
	cmd := exec.Command("bash", fmt.Sprint("install/", script))
	byteoutput, err := cmd.Output()
	if err != nil {
		return err
	}
	output = string(byteoutput)

	// display script output to user
	fmt.Print(output)
	return nil

}

// Used for commands that are simply handled by the shell function
func passToShell(data *CmdArgs) error {
	c := data.cmd[0]
	command := string(c[1:])

	switch command {
	case "]", "[", "..", "-", "hist":
		fmt.Println(command)
	default:
		return errors.New(fmt.Sprintf("Tried to pass command to shell, but '%s' is not a valid command for the shell function.", command))
	}

	return nil
}
