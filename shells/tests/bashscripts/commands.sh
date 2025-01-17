
test_create_project_tree() {
    mkdir something
    cd something
    mkdir else
    cd else
    mkdir another
    cd another
    cd ../../..
}

available_commands=("]" "[" "-hist" ".." "-")

test_cmd() {
    local selected=$(printf "%s\n" "${available_commands[@]}" | fzf)
    if [[ -n "$selected" ]]; then
        ft__execute "$selected"
    fi
}
