

# based on user's os, determine correct directory for exe to save to
# create fastTravel directory at TARGET_DIR and save exe to this location
# add main.sh contents to users shell rc 

# TODO
# add step to test executable after install by running `version` command (ps, add `version` command) 
# add manual install option where user can choose what shell(s) to install fastTravel to



SCRIPT_PATH="./main.sh"



function find_target_dir() {

    local os=$(uname)
    case "$os" in
        Linux*)
            echo "/usr/local/bin"
            ;;
        Darwin*)
            echo "/usr/local/bin"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            echo "%USERPROFILE%/AppData/Local"  
            ;;
        *)
            echo "Error! Unsupported operating system"
            exit 1
            ;;
    esac

}


TARGET_DIR=$(find_target_dir)


exe_install() {
    
    ft_dir="$TARGET_DIR/fastTravel"

    mkdir -p "$ft_dir"

    mv ./fastTravel.exe "$ft_dir/" 

    echo "$ft_dir/fastTravel.exe"
}

FT_EXE_PATH=$(exe_install)



# USER_SHELL=""



function bash_install() {
    case "$SHELL" in
        *bash*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.bashrc
            echo ". $SCRIPT_PATH" >> ~/.bashrc
            ;;
        *zsh*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.zshrc
            echo ". $SCRIPT_PATH" >> ~/.zshrc
            ;;
        *fish*)
            echo "set -x FT_EXE_PATH \"$FT_EXE_PATH\"" >> ~/.config/fish/config.fish
            echo "source $SCRIPT_PATH" >> ~/.config/fish/config.fish
            ;;
        *csh* | *tcsh*)
            echo "setenv FT_EXE_PATH \"$FT_EXE_PATH\"" >> ~/.cshrc
            echo "source $SCRIPT_PATH" >> ~/.cshrc
            ;;
        *ksh* | *sh*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.kshrc
            echo ". $SCRIPT_PATH" >> ~/.kshrc
            ;;
        *powershell*)
            # PowerShell profile path varies depending on the version
            if [ -f "$%USERPROFILE%/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1" ]; then
                echo "\$FT_EXE_PATH = \"$FT_EXE_PATH\"" >> "$%USERPROFILE%/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
                echo ". '$SCRIPT_PATH'" >> "$%USERPROFILE%/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
            elif [ -f "$%USERPROFILE%/Documents/PowerShell/Microsoft.PowerShell_profile.ps1" ]; then
                echo "\$FT_EXE_PATH = \"$FT_EXE_PATH\"" >> "$%USERPROFILE%/Documents/PowerShell/Microsoft.PowerShell_profile.ps1"
                echo ". '$SCRIPT_PATH'" >> "$%USERPROFILE%/Documents/PowerShell/Microsoft.PowerShell_profile.ps1"
            else
                echo "echo \$FT_EXE_PATH = \"$FT_EXE_PATH\" >> PowerShell profile file"
                echo "echo . '$SCRIPT_PATH' >> PowerShell profile file"
            fi
            ;;
        *)
            echo "Unrecognized shell. Please add the script to your shell's configuration file manually."
            exit 1
            ;;
    esac
}



bash_install




