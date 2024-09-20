#!/usr/bin/env bash

upperStack=()

pushup() {
    local _path="$1"
    upperStack+=("$_path")
}

popup() {
    unset 'upperStack[-1]'
    upperStack=(${upperStack[@]})
}

ft() {

    temp_output=$(mktemp)
    
    "$FT_EXE_PATH" "$@" | tee "$temp_output"

    output="$(tail -n 1 "$temp_output")"

    if [[ -d "$output" || "$output" == ".." || "$output" == "-" ]]; then 
        upperStack=()
        cd "$output";
    elif [[ "$output" == "]" ]]; then    
        if [ ${#upperStack[@]} -eq 0 ]; then
            echo Already at head of history stack.
            rm "$temp_output"
            return 1
        
        fi
        local p="${upperStack[-1]}"
        popup
        cd "$p"

    elif [[ "$output" == "[" ]]; then
        local lowerStackLen=$(dirs -v | awk '{print $1}' | sort -n | tail -1)
        if [ "$lowerStackLen" -eq 0 ]; then
            echo Already at tail of history stack.
            rm "$temp_output"
            return 1
        fi

        local p=$(pwd)
        pushup "${p}"
        popd > /dev/null
    
    fi

    rm "$temp_output"

}



