#!/usr/bin/env bash

ft__check_fzf() {
    if ! command -v fzf &>/dev/null; then
        echo "Error: fzf is needed to use this feature. https://github.com/junegunn/fzf"
        return 1
    fi
}

ft__check_tree() {
    if ! command -v tree &>/dev/null; then
        echo "Error: tree is needed to use this feature. You can install it using your package manager:"
        echo "  - Ubuntu/Debian: sudo apt install tree"
        echo "  - Fedora: sudo dnf install tree"
        echo "  - Arch: sudo pacman -S tree"
        echo "  - macOS: brew install tree"
        return 1
    fi
}

ft__upperStack=()

ft__pushup() {
    local path="$1"
    ft__upperStack+=("$path")
}

ft__popup() {
    unset 'ft__upperStack[-1]'
    ft__upperStack=(${ft__upperStack[@]})
}



ft() {

    local temp_output=$(mktemp)

    "$FT_EXE_PATH" "$@" | tee "$temp_output"
    local output="$(tail -n 1 "$temp_output")"
    rm "$temp_output"

    if [[ -d "$output" || "$output" == ".." || "$output" == "-" ]]; then 
        
        ft__upperStack=()
        pushd "$output" > /dev/null

    elif [[ "$output" == "]" ]]; then    
        
        if [ ${#ft__upperStack[@]} -eq 0 ]; then
            echo Already at head of history stack.
            return 1 
        fi

        local p="${ft__upperStack[-1]}"
        ft__popup
        pushd "$p" > /dev/null

    elif [[ "$output" == "[" ]]; then
        
        local lowerStackLen=$(dirs -v | awk '{print $1}' | sort -n | tail -1)
        
        if [ "$lowerStackLen" -eq 0 ]; then
            echo Already at tail of history stack.
            return 1
        fi

        local p=$(pwd)
        ft__pushup "${p}"
        popd > /dev/null
    
    elif [[ "$output" == "hist" ]]; then
        ft__check_fzf
        ft__check_tree
        
        local lowerStack=($(dirs -v | awk '{print $2}'))
        local historyStack=("${ft__upperStack[@]}" "${lowerStack[@]}")

        printf "%s\n" "${historyStack[@]}" | fzf \
            --tac --header "Currently at $(pwd)" --preview 'tree {}' | while read -r dir; do eval ft "$dir"; done

    fi    

}



