package main

import (
    "fmt"
    "os"
)


func passCmd(args []string) []string {
    
    if (len(args)) >= 2 {
        return args[1:]
    } else {
        fmt.Println("No command provided. Usage: ftravel <command> <path>)
        os.Exit(1)
    }

}

func changeDirectory(path string) error {
    
    err := os.Chdir(path)
    if err != nil {
        fmt.Printf("Fast travel failed! %v", err)
        os.Exit(1)
    }

    fmt.Println("moved to ", path)

}


func main() {
    
    availCmds := map[string]func(path string) {
        "to": changeDirectory,
    }
    
    inputCommand := passCmd(os.Args)
    action := input_command[0]
    path := input_command[1]

    
    exeCmd, ok := availCmds[action]
    if !ok {
        fmt.Println("Invalid command, use 'help' for available commands.")
        os.Exit(1)
    }

    exeCmd(path)

}



