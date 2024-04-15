


SCRIPT_PATH="ftmain.sh"



find_target_dir() {

    local os=$(uname)
    case "$os" in
        Linux*)
            echo "/usr/local/bin"
            ;;
        Darwin*)
            echo "/usr/local/bin"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            echo "$HOME/AppData/Local"  
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

    sudo mkdir -p "$ft_dir"

    sudo mv ./fastTravel.exe "$ft_dir/" 

    echo "$ft_dir/fastTravel.exe"
}

FT_EXE_PATH=$(exe_install)

export PATH="$PATH:$FT_EXE_PATH"

# USER_SHELL=""



bash_install() {
    case "$SHELL" in
        *bash*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.bashrc
            echo ". ~/$SCRIPT_PATH" >> ~/.bashrc
            ;;
        *zsh*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.zshrc
            echo ". ~/$SCRIPT_PATH" >> ~/.zshrc
            ;;
        *fish*)
            echo "set -x FT_EXE_PATH \"$FT_EXE_PATH\"" >> ~/.config/fish/config.fish
            echo "source ~/$SCRIPT_PATH" >> ~/.config/fish/config.fish
            ;;
        *csh* | *tcsh*)
            echo "setenv FT_EXE_PATH \"$FT_EXE_PATH\"" >> ~/.cshrc
            echo "source ~/$SCRIPT_PATH" >> ~/.cshrc
            ;;
        *ksh* | *sh*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.kshrc
            echo ". ~/$SCRIPT_PATH" >> ~/.kshrc
            ;;
        *powershell*)
            # PowerShell profile path varies depending on the version
            if [ -f "$HOME/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1" ]; then
                echo "\$FT_EXE_PATH = \"$FT_EXE_PATH\"" >> "$HOME/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
                echo ". '~/$SCRIPT_PATH'" >> "$HOME/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
            elif [ -f "$HOME/Documents/PowerShell/Microsoft.PowerShell_profile.ps1" ]; then
                echo "\$FT_EXE_PATH = \"$FT_EXE_PATH\"" >> "$HOME/Documents/PowerShell/Microsoft.PowerShell_profile.ps1"
                echo ". '~/$SCRIPT_PATH'" >> "$HOME/Documents/PowerShell/Microsoft.PowerShell_profile.ps1"
            else
                New-Item "$HOME/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1" 
                echo "\$FT_EXE_PATH = \"$FT_EXE_PATH\"" >> "$HOME/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
                echo ". '~/$SCRIPT_PATH'" >> "$HOME/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1"
            fi
            ;;
        *)
            echo "Unrecognized shell. Please add the script to your shell's configuration file manually."
            exit 1
            ;;
    esac

    sudo cp $SCRIPT_PATH ~/

}



bash_install




