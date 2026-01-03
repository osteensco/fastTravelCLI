package ft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)


var QueryCascade = map[string]struct {
	FuncMap func(allpaths map[string]string, path string) (string, error)
}{
	// TODO: make these keys enums?
	"bookmark": {evalBookmark},
	"CDPATH": {evalCDPATH},
	"relative": {evalRelative},
	// "projectBookmark": {evalProjectBookmark},
	// "zoxide": {evalZoxide},
}


// TODO
//  - check for SPECIFIC errors in directory checks (ex: *PathError)
//  - ensure other types of errors are caught and returned early

// evaluates a bookmark query and returns an absolute path
func evalBookmark(allPaths map[string]string, path string) (string, error) {
	var key string
	provided_string := path

	if strings.Contains(provided_string, "/") {

		path_array := strings.Split(provided_string, "/")
		eval_array := make([]string, len(path_array))

		// key evaluation if the first string before "/" delimeter is key
		key = path_array[0]
		if p, ok := allPaths[key]; ok {
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

		return "", fmt.Errorf(InvalidDirectoryMsg, provided_string, path)

	} else {

		key = provided_string
		// handles key lookup
		p, ok := allPaths[key]
		
		if ok {
			return fmt.Sprintf("%s", p), nil
		}
	}
	
	return "", nil

}

// evaluates a CDPATH query and returns an absolute path
func evalCDPATH(allpaths map[string]string, path string) (string, error) {
	cdpath := os.Getenv("CDPATH")
	if cdpath != "" {
		cdpaths := strings.Split(cdpath, ":")
		for _, cdp := range cdpaths {
			cdPathResult := filepath.Join(cdp, path)
			dir, err := os.Stat(cdPathResult)
			if err == nil && dir.IsDir() {
				return fmt.Sprintf("%s", cdPathResult), nil
			}
		}
	}
	return "", nil
}

// evaluates a relative path query and returns an absolute path
func evalRelative(allpaths map[string]string, path string) (string, error) {
	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			// in this case key is assumed to be a relative path
			p, err := filepath.Abs(path)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s", p), nil
		}
	}
	return "", nil
}

// TODO
// func evalProjectBookmark(allpaths *[]string, path *string) error {
// 	return nil
// }
//
// func evalZoxide(allpaths *[]string, path *string) error {
// 	return nil
// }


/*
Takes a path query and returns it's absolute path or an error.
*/
func evalPath(data *CmdAPI, path *string) (string, error) {

	var res string
	var err error

	for _, query := range data.settings.QueryOrder {
		res, err = QueryCascade[query].FuncMap(data.allPaths, *path)
		if err != nil {
			return res, err
		}
		if res != "" {
			return res, nil
		}
	}

	return "", errors.New(fmt.Sprintf(UnrecognizedKeyMsg, *path))


	// // ----------- bookmarks and key evals -----------
	// var key string
	// provided_string := *path
	//
	// if strings.Contains(provided_string, "/") {
	//
	// 	path_array := strings.Split(provided_string, "/")
	// 	eval_array := make([]string, len(path_array))
	//
	// 	// key evaluation if the first string before "/" delimeter is key
	// 	key = path_array[0]
	// 	if p, ok := data.allPaths[key]; ok {
	// 		eval_array[0] = p
	// 	} else {
	// 		eval_array[0] = key
	// 	}
	// 	for i, str := range path_array {
	// 		if i != 0 {
	// 			eval_array[i] = str
	// 		}
	// 	}
	//
	// 	// handles evaluated path
	// 	path := strings.Join(eval_array, "/")
	// 	dir, err := os.Stat(path)
	// 	if err == nil {
	// 		if dir.IsDir() {
	// 			return fmt.Sprintf("%s", path), nil
	// 		}
	// 	}
	//
	// 	return "", errors.New(fmt.Sprintf(InvalidDirectoryMsg, provided_string, path))
	//
	// } else {
	//
	// 	key = provided_string
	// 	// handles key lookup
	// 	p, ok := data.allPaths[key]
	//
	// 	// ----------------------------------------------
	// 	if !ok {
	// 		// handles releative directory in CWD
	// 		// -----------relative paths---------------------
	// 		dir, err := os.Stat(key)
	// 		if err == nil {
	// 			if dir.IsDir() {
	// 				// in this case key is assumed to be a relative path
	// 				p, err = filepath.Abs(key)
	// 				if err != nil {
	// 					return "", err
	// 				}
	// 				return fmt.Sprintf("%s", p), nil
	// 			}
	// 		}
	// 		// ----------------------------------------------
	//
	// 		// ----------CDPATH------------------------------
	// 		cdpath := os.Getenv("CDPATH")
	// 		if cdpath != "" {
	// 			cdpaths := strings.Split(cdpath, ":")
	// 			for _, path := range cdpaths {
	// 				cdPathResult := filepath.Join(path, key)
	// 				dir, err := os.Stat(cdPathResult)
	// 				if err == nil && dir.IsDir() {
	// 					return fmt.Sprintf("%s", cdPathResult), nil
	// 				}
	// 			}
	// 		}
	// 		// ----------------------------------------------
	//
	// 		return "", errors.New(fmt.Sprintf(UnrecognizedKeyMsg, key))
	//
	// 	}
	//
	// 	return fmt.Sprintf("%s", p), nil
	//
	// }
}
