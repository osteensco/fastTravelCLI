
maketree() {
    mkdir something
    pushd something
    mkdir else
    pushd else
    mkdir another
    pushd another
    pushd ../../..
}


available_commands=("[" "]" ".." "-" "hist")

# Spoof fastTravelCLI's executable output
export FT_EXE_PATH="/testspace/exe.sh"

testcmd() {
    local selected=$(printf "%s\n" "${available_commands[@]}" | fzf)
    if [[ -n "$selected" ]]; then
        ft "$selected"
    fi
}


