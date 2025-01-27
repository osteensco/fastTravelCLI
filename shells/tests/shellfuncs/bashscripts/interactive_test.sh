#!/bin/bash


# This test script is designed to test the ftmain.sh file independent of fastTravelCLI's binary.
# Use by calling the simtest function.
# This script will build simulations to run through automatically.
# setup.sh and ftmain.sh need to be sourced prior to execution of these simulations. This step is added to the .bashrc in Docker.
# To manually test specific commands, utilize the testcmd function found in setup.sh.



# Spoof fastTravelCLI's executable output
export FT_EXE_PATH="/testspace/exe.sh"

commands=("[" "]" ".." "-" "hist")



# Use this function for testing commands manually.
testcmd() {
    local selected=$(printf "%s\n" "${commands[@]}" | fzf)
    if [[ -n "$selected" ]]; then
        ft "$selected"
    fi
}



simtest() {
    # User can set number of simulations or use default of 15
    if [[ $# -ne 1 ]]; then
        sims=15
    else
        if ! [[ "$1" =~ ^[0-9]+$ ]]; then
            echo "Argument must be a number! Arg: $1" >&2
            exit 1
        fi
        sims=$1
    fi

    # Generate a random simulation
    # Setup script will define the commands
    generateSim() {
        local totalcmds=$1
        # Start the simulation more towards the middle of the history stack
        local sim=("[" "[" "[" "[" "[" "[")
        for ((i = 0; i < totalcmds; i++)); do
            local idx=$(( RANDOM % ${#commands[@]} ))
            sim+=("${commands[idx]}")
        done
        echo "${sim[@]}"
    }

    # Create the simulation
    simulation=($(generateSim "$sims"))

    # Process the commands
    for cmd in "${simulation[@]}"; do
        echo -e "--- History stack: \n $(ft__phist)\n-----"
        echo -e "--- Currently at directory: $(pwd)\n--->"
        
        ft "$cmd"

        echo "<---"

    done
}
