package ft

import (
	"bufio"
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
		return false, fmt.Errorf("%v is not a valid response, please type y/n", s)
	}
}

func printMap(hashmap map[string]string, formatstr string) {

	keys := make([]string, 0, len(hashmap))
	for k := range hashmap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i := range keys {

		fmt.Printf(formatstr, keys[i], hashmap[keys[i]])

	}

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

	// exit early if no piped input is detected (stdin is terminal)
	if fileinfo.Mode()&os.ModeCharDevice != 0 {
		return nil
	}

	// read from stdin (piped input)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		pipedArgs := splitWords(line) // split into individual arguments
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
			// the first arg should always correspond to the shell's function call (default 'ft')
			// we don't explicitly check the function call since a user could potentially use an alias or wrap it in another call
			continue
		}
		// command and flags will have the '-' prefix
		// acts as a 'permissive' checker
		if v[0] == '-' {
			v = strings.ToLower(v)
			// handle flag parsing
			switch v {
			case "-y":
				cmd.Flags.Y = true
				// while help is a stand alone command, it can be passed as a flag to get detailed help on other commands.
			case "-h", "-help":
				cmd.Flags.H = true
			default:
				// if not a valid flag, we assume it's a command
				// validity of command is checked elsewhere in the passCmd function
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
	// handle -help as a command
	if cmd.Flags.H && cmd.Cmd == "" {
		cmd.Flags.H = false
		cmd.Cmd = "-help"
	}

	return cmd
}

func CreateCmdHelpDoc(DescSlice []string) string {
	cmd := DescSlice[0]
	desc := DescSlice[1]
	return fmt.Sprintf("\nUsage: %s\n\n%s", cmd, desc)
}

func CreateHelpOutput() string {

	output := "\nUsage:\n"
	for _, val := range CmdDesc {
		for k, v := range val {
			output += fmt.Sprintf(HelpLineStrFormat, k, v)
		}
	}
	output += HelpDescExamples
	output += "\nFor more detailed help for a specific command, add -h or -help to that command (e.g. ft -set -h)\n"

	return output
}

func DisplayDetailedHelp(action string) string {
	help, ok := DetailedCmdDescMapping[action]
	if !ok {
		panic(fmt.Sprintf("Command does not have a DetailedCmdDescMapping - '%s'", action))
	}
	return help
}

func EnsureLength(actual int, expected int) {
	if actual != expected {
		panic(fmt.Sprintf("DetailedCmdDescriptions length expected to be %v, actual %v", expected, actual))
	}
}

func FindUsageMaxLen(usages map[string]string) int {
	maxLen := 0
	for _, u := range usages {
		if len(u) > maxLen {
			maxLen = len(u)
		}
	}

	return maxLen + 2
}
