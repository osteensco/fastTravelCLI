#!/usr/bin/env bash

ft__check_fzf() {
    if ! command -v fzf &>/dev/null; then
        echo "Error: fzf is needed to use this feature. You can install it using your package manager:"
        echo "  - Ubuntu/Debian: sudo apt install fzf"
        echo "  - Fedora: sudo dnf install fzf"
        echo "  - Arch: sudo pacman -S fzf"
        echo "  - macOS: brew install fzf"
        return 1
    fi
}

ft__upperStack=()

ft__pushup() {
    local _path="$1"
    ft__upperStack+=("$_path")
}

ft__popup() {
    unset 'ft__upperStack[-1]'
    ft__upperStack=(${ft_upperStack[@]})
}

ft__capture() {
    local temp_output=$(mktemp)
    
    "$FT_EXE_PATH" "$@" | tee "$temp_output"
    local output="$(tail -n 1 "$temp_output")"

    echo "$output"
}

ft__execute() {
    local output="$1"
    echo "$output"

    if [[ -d "$output" || "$output" == ".." || "$output" == "-" ]]; then 
        
        ft__upperStack=()
        pushd "$output" > /dev/null

    elif [[ "$output" == "]" ]]; then    
        
        if [ ${#ft__upperStack[@]} -eq 0 ]; then
            echo Already at head of history stack.
            rm "$temp_output"
            return 1
        
        fi

        local p="${ft__upperStack[-1]}"
        ft__popup
        pushd "$p" > /dev/null

    elif [[ "$output" == "[" ]]; then
        
        local lowerStackLen=$(dirs -v | awk '{print $1}' | sort -n | tail -1)
        
        if [ "$lowerStackLen" -eq 0 ]; then
            echo Already at tail of history stack.
            rm "$temp_output"
            return 1
        fi

        local p=$(pwd)
        ft__pushup "${p}"
        popd > /dev/null
    
    elif [[ "$output" == "-hist" ]]; then
        # TODO: test me!!
        ft__check_fzf
        
        local navigation="up;"
        for i in "${ft__upperStack[@]}"; do
            navigation="$navigation;"
            navigation="$navigation down"
        done

        local lowerStack=($(dirs -v | awk '{print $2}'))
        local historyStack=("${ft__upperStack[@]}" "${lowerStack[@]}")

        local selected=$(printf "%s\n" "${historyStack[@]}" | fzf \
            --bind "start:execute-silent($navigation)")
        
        if [[ -n "$selected" ]]; then
            pushd "$selected" > /dev/null || echo "Could not find directory $selected"
        fi

    fi
    
    if [[ -f "$temp_output" ]]; then
        rm "$temp_output"
    fi

}


ft() {
    local output=$(ft__capture "$@")
    ft__execute "$output"
}

