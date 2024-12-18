package ft

import (
	"bufio"
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
		return false, errors.New(fmt.Sprintf("%v is not a valid response, please type y/n", s))
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

func PipeArgs(args *[]string) error {
	fileinfo, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	// Exit early if no piped input is detected
	if fileinfo.Mode()&os.ModeCharDevice != 0 {
		return nil
	}

	// read from stdin
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		// append to args
		*args = append(*args, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil

}
