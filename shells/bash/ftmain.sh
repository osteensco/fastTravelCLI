#!/usr/bin/env bash

upperStack=()

pushup() {
    local path=$1
    upperStack+=("$path")    
}

popup() {
    local path="${upperStack[-1]}"
    unset 'upperStack[-1]'
    echo "$path"
}

ft() {

    temp_output=$(mktemp)
    
    "$FT_EXE_PATH" "$@" | tee "$temp_output"

    output="$(tail -n 1 "$temp_output")"

    if [ -d "$output" ]; then 
        cd "$output";
    elif [[ "$output" == ">" ]]; then    
        if [ ${#upperStack[@]} -eq 0 ]; then
            echo Already at head of history stack.
            return 1
        fi
        p=$(popup)
        cd "$p"
    elif [[ "$output" == "<" ]]; then
        p=$(dirs +0)
        pushup "$p"
        popd > /dev/null
    fi

    rm "$temp_output"

}



