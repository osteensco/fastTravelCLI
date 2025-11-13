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

ft__phist() {
    # print fastTravelCLI's history stack
    local lowerStack=($(dirs -v | awk '{print $2}'))
    local historyStack=("${ft__upperStack[@]}" "${lowerStack[@]}")

    printf "%s\n" "${historyStack[@]}" 
}

ft__handle_fzf_dirs() {
    # handle more complex fzf integrated actions
    local output_array=(${=1})
    local fzf_cmd="${output_array[1]}"
    local dir=""
    if (( ${#output_array[@]} == 1 )); then
        dir="."
    elif (( ${#output_array[@]} > 1)); then
        dir="${output_array[2]}"
    fi
    # if [[ ! -d "$dir" ]]; then
    #     echo "Error: '$dir' is not a valid directory"
    #     return 1
    # fi
    
    case "$fzf_cmd" in 
        fzfc)
            ft__check_fzf
            ft__check_tree
            # fuzzy find current level
            eval ft "$(ls "$dir" -d */ | fzf --tac --header " fastTravelCLI " --border-label " fastTravelCLI " --preview 'tree -L 1 {}')"
            ;;

        fzfa)
            ft__check_fzf
            ft__check_tree
            # subdir exclusion list
            EXCLUDES=(.venv venv .env env node_modules __pycache__ .git .svn .hg .idea .vscode .DS_Store dist build out target coverage .cache tmp temp)
            # build find command
            find_cmd=(find "$dir" -type d)
            for e in "${EXCLUDES[@]}"; do
                find_cmd+=(-not -path "*/$e/*")
            done
            # build tree command
            OIFS=$IFS
            IGNORE_PATTERN=$(IFS='|'; echo "${EXCLUDES[*]}")
            IFS=$OIFS
            # fuzzy find all lower subdirs
            eval ft "$("${find_cmd[@]}" | fzf --tac --header " fastTravelCLI " --border-label " fastTravelCLI " --preview "tree -L 1 -I \"$IGNORE_PATTERN\" {}")"
            ;;
    esac
}

ft__handle_fzf_bookmarks() {
    ft__check_fzf
    ft__check_tree
    # TODO generate a different view for no bookmarks?
    # fuzzy find bookmarks
    eval ft -ls | fzf --tac --header " fastTravelCLI - Bookmarks" --border-label " fastTravelCLI - Bookmarks" --delimiter=": " --preview 'tree -L 1 {2}' | cut -d':' -f1 | ft
}

ft__handle_fzf_hist() {
    ft__check_fzf
    ft__check_tree
    # fuzzy find session history
    eval ft "$(ft__phist | fzf --tac --header " fastTravelCLI - History" --border-label " fastTravelCLI - History" --preview 'tree -L 1 {}')"
}



ft() {

    
    # output cache
    local temp_output=$(mktemp)
    
    # run fastTravelCLI's binary and capture output
    "$FT_EXE_PATH" "$@" | tee "$temp_output"
    # truncate any irrelevant output
    local output="$(tail -n 1 "$temp_output")"
    rm "$temp_output"

    # handle base functionality (simple navigation)
    if [[ -d "$output" || "$output" == ".." || "$output" == "-" ]]; then 
        
        ft__upperStack=()
        pushd "$output" > /dev/null

    # handle history navigation forwards
    elif [[ "$output" == "]" ]]; then    
        
        if [ ${#ft__upperStack[@]} -eq 0 ]; then
            echo Already at head of history stack.
            return 0 
        fi

        local p="${ft__upperStack[-1]}"
        ft__popup
        pushd "$p" > /dev/null

    # handle history navigation backwards
    elif [[ "$output" == "[" ]]; then
        
        local lowerStackLen=$(dirs -v | awk '{print $1}' | sort -n | tail -1)
        
        if [ "$lowerStackLen" -eq 0 ]; then
            echo Already at tail of history stack.
            return 0
        fi

        local p=$(pwd)
        ft__pushup "${p}"
        popd > /dev/null
   
    # handle fzf for history stack
    elif [[ "$output" == "hist" ]]; then
        ft__handle_fzf_hist 
    
    # handle fzf for bookmarks
    elif [[ "$output" == "fzf" ]]; then
        ft__handle_fzf_bookmarks 

    # handle fzf for dirs in current project
    else
        ft__handle_fzf_dirs "$output"
    fi    
}



