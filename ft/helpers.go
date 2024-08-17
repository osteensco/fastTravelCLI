package ft

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)



func verifyInput(s string) (bool, error) {
    switch strings.ToLower(s) {
        case "y":
            return true, nil
        case "n":
            return false, nil
        default:
            return false, errors.New(fmt.Sprintf("%v is not a valid response, please type y/n",s))
    } 
}

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

func SanitizeDir(path string) string {
	distro := os.Getenv("WSL_DISTRO_NAME")
	if len(distro) == 0 {

    	prefix := "/mnt/"
        drive := strings.Split(path, ":")[0]
        path = strings.Replace(path, drive, strings.ToLower(drive), 1)
        path = strings.Replace(path, ":", "", 1)
    	path = strings.Replace(path, "\\", "/", -1)
    	path = prefix + path

	}
    return path
}


func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}


