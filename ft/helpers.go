package ft

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

func verifyInput(s string, f bool) (bool, error) {
	if f {
		return true, nil
	}
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

// Parses os.Args and organizes command into a struct
func ParseArgs(args *[]string) *Cmd {
	cmd := NewCmd(args)
	for i, v := range *args {
		if i == 0 {
			// the first arg should always correspond to the shell's function call
			// we don't explicitly check the function call (default 'ft') since a user could potentially use an alias or wrap it in another call
			continue
		}
		// command and flags will have the '-' prefix
		// acts as a 'permissive' checker
		if v[0] == '-' {
			v = strings.ToLower(v)
			switch v {
			case "-y":
				cmd.Flags.Y = true
			default:
				// if not a valid flag, we assume it's a command
				// validity of command is checked elsewhere
				// first command found is the one identified for use
				if cmd.Cmd != "" {
					continue
				}
				cmd.Cmd = v
			}
			continue
		} else {
			// handle special navigation commands
			switch v {
			case "]", "[", "..", "-":
				if cmd.Cmd == "" {
					cmd.Cmd = v
				}
			default:

				cmd.Args = append(cmd.Args, v)
			}
		}
	}

	return cmd
}
