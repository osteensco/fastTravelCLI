#!/bin/bash


# This test script is designed to test the ftmain.sh file independent of fastTravelCLI's binary.
# This script will run through predefined set of commands to simulate a user journey.
# This is designed to run automatically with no interactivity in a docker container.



# Source required scripts
source ~/.fzf.bash
source setup.sh
source ftmain.sh



# TODO
#  - change this hashmap into two arrays so that order is maintained.
#  - add check for 'already top of stack' and 'already bottom of stack'

# Test commands and their expected outputs
test_cmds=(
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft [" 
    "ft ]" 
    "ft ]" 
    "ft ]" 
    "ft .." 
    "ft -"
    "ft ]"
)

test_expected=(
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[" 
    "[\nAlready at tail of history stack." 
    "]" 
    "]" 
    "]" 
    ".." 
    "-" 
    "]\nAlready at head of history stack."
)

# Variable to hold overall test result
all_tests_passed=true

# Run commands, capture output, and compare to expected
i=0
for command in "${test_cmds[@]}"; do

    # Create a temporary file to capture output
    tempfile=$(mktemp)
    # Run the command in the current shell, redirecting output to the tempfile
    eval "$command" > "$tempfile" 2>&1
    # Read the captured output into a variable
    output=$(<"$tempfile")
    # Clean up the temp file
    rm "$tempfile"

    # Compare output to expected
    expected=$(echo -e "${test_expected[$i]}")
    if [[ "$output" == "$expected" ]]; then
    echo "PASS: '$command' output matched expected."
    else
    echo "FAIL: '$command' output did not match expected."
    echo "  Expected: $expected"
    echo "  Got: $output"
    all_tests_passed=false
    fi
    i=$((i+1))
done

# Simulate fzf for the `ft hist` command
export FZF_DEFAULT_OPTS="--filter=testspace" # Automatically select the first match
hist_output=$(ft hist 2>&1 | head -n 2)
expected_hist_output=$(echo -e "hist\n/testspace")

# Compare fzf-based output
if [[ "$hist_output" == "$expected_hist_output" ]]; then
    echo "PASS: 'ft hist' output matched expected."
else
    echo "FAIL: 'ft hist' output did not match expected."
    echo "  Expected: $expected_hist_output"
    echo "  Got: $hist_output"
    all_tests_passed=false
fi

# Overall test result
if $all_tests_passed; then
    echo "All tests passed!"
else
    echo "Some tests failed."
    exit 1
fi

