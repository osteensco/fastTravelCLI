#!/usr/bin/env bash

check_fzf() {
    if ! command -v fzf &>/dev/null; then
        echo "Error: fzf is needed to use this feature. You can install it using your package manager:"
        echo "  - Ubuntu/Debian: sudo apt install fzf"
        echo "  - Fedora: sudo dnf install fzf"
        echo "  - Arch: sudo pacman -S fzf"
        echo "  - macOS: brew install fzf"
        return 1
    fi
}

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
    
    elif [[ "$output" == "-hist" ]]; then
        # TODO: test me!!
        check_fzf
        
        local navigation="up;"
        for i in "${upperStack[@]}"; do
            navigation="$navigation;"
            navigation="$navigation down"
        done

        local lowerStack=($(dirs -v | awk '{print $2}'))
        local historyStack=("${upperStack[@]}" "${lowerStack[@]}")

        local selected=$(printf "%s\n" "${historyStack[@]}" | fzf \
            --bind "start:execute-silent($navigation)")
        
        if [[ -n "$selected" ]]; then
            cd "$selected" || echo "Could not find directory $selected"
        fi

    fi

    rm "$temp_output"

}



