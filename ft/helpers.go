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

// Split a string by spaces similar to strings.Split(), but treating everything within quotes as one word.
func splitWords(line string) []string {
	var result []string
	var current strings.Builder
	inquote := false
	var quoteChar rune

	for i, r := range line {
		switch r {
		case '"', '\'':
			if inquote {
				if r == quoteChar {
					inquote = false
				} else {
					current.WriteRune(r)
				}
			} else {
				inquote = true
				quoteChar = r
			}
		case ' ':
			if inquote {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}

		if i == len(line)-1 && current.Len() > 0 {
			result = append(result, current.String())
		}
	}

	return result
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
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {

		line := scanner.Text()
		pipedArgs := splitWords(line)

		// append to args
		*args = append(*args, pipedArgs...)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil

}
