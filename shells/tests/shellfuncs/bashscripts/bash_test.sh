#!/bin/bash


# This test script is designed to test the ftmain.sh file independent of fastTravelCLI's binary.
# This script will run through predefined set of commands to simulate a user journey.
# This is designed to run automatically with no interactivity in a docker container.



# Source dependency scripts
source ~/.fzf.bash
source ./maketree.sh
source ./shells/bash/ftmain.sh



# Spoof fastTravelCLI's executable output
export FT_EXE_PATH="/testspace/tests/bashscripts/exe.sh"



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

# Loop over scripted simulation
i=0
for command in "${test_cmds[@]}"; do
    
    # Capture output
    tempfile=$(mktemp)
    eval "$command" > "$tempfile" 2>&1
    output=$(<"$tempfile")
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
hist_output=$(ft hist 2>&1 | head -n 2) # stdout will capture fzf directory check, so we only want first two lines since the rest is hidden to the user
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



if $all_tests_passed; then
    echo "Bash tests completed with no failures."
else
    echo "Bash tests completed with failures."
    exit 1
fi

