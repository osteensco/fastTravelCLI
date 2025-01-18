
maketree() {
    mkdir something
    cd something
    mkdir else
    cd else
    mkdir another
    cd another
    cd ../../..
}

maketree

available_commands=("]" "[" "hist" ".." "-")

testcmd() {
    local selected=$(printf "%s\n" "${available_commands[@]}" | fzf)
    if [[ -n "$selected" ]]; then
        ft__execute "$selected"
    fi
}


