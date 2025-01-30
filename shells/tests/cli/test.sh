#!/bin/bash

# Wrapper script to run all cli tests sequentially.




set -eo pipefail

declare -A TESTS=(
  ["./tests/bash/bash_test.sh"]="/bin/bash"
  # ["./tests/zsh/zsh_test.sh"]="/bin/zsh"
)

all_tests_passed=true

for TEST in "${!TESTS[@]}"; do
    SHELL="${TESTS[$TEST]}"
    echo "Running cli test $TEST with $SHELL..."

    if ! command -v "$SHELL" >/dev/null 2>&1; then
        echo "ERROR: Shell $SHELL not found!"
        exit 1
    fi

    if "$SHELL" -c "set -e; $TEST"; then
        echo -e "PASSED: $TEST"
    else
        echo -e "FAILED: $TEST (exit code $?)"
        all_tests_passed=false
    fi
done


if $all_tests_passed; then
    echo "All cli tests passed!"
else
    echo "Some cli tests failed."
    exit 1
fi


