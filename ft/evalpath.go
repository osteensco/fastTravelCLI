package ft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TODO:
//  - move evalPath to this file
//  - setup some sort of map for individual functions identified
//  - change logic in evalPath to read in settings and execute given order
//  - move tests to new tests file
//  - add tests for various settings and reading in settings

/*
Takes a key or relative path and returns it's absolute path or an error.
*/
func evalPath(data *CmdAPI, path *string) (string, error) {

	// TODO
	//  - check for SPECIFIC errors in directory checks (ex: *PathError)
	//  - ensure other types of errors are caught and returned early

	// TODO: break out into individual functions so settings can determine order
	//	- current: bookmark, CDPATH, relative path
	//	- future: project bookmark, zoxide query

	var key string
	provided_string := *path

	if strings.Contains(provided_string, "/") {

		// ----------- bookmarks and key evals -----------
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

		// handles evaluated path
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

		// ----------------------------------------------
		if !ok {
			// handles releative directory in CWD
			// -----------relative paths---------------------
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
			// ----------------------------------------------

			// ----------CDPATH------------------------------
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
			// ----------------------------------------------

			return "", errors.New(fmt.Sprintf(UnrecognizedKeyMsg, key))

		}

		return fmt.Sprintf("%s", p), nil

	}
}
