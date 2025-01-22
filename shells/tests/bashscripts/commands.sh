
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

# Spoof fastTravelCLI's executable output
export FT_EXE_PATH="/testspace/exe.sh"

testcmd() {
    local selected=$(printf "%s\n" "${available_commands[@]}" | fzf)
    if [[ -n "$selected" ]]; then
        ft "$selected"
    fi
}


