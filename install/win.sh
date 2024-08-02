



SCRIPT_PATH="ftmain.sh"

TARGET_DIR="$HOME/AppData/Local"


exe_install() {
    
    ft_dir="$TARGET_DIR/fastTravel"

    sudo mkdir -p "$ft_dir"

    sudo mv ./fastTravel.exe "$ft_dir/" 

    echo "$ft_dir/fastTravel.exe"
}

FT_EXE_PATH=$(exe_install)

export PATH="$PATH:$FT_EXE_PATH"

USER_SHELL="echo $SHELL"
if command -v pwsh &> /dev/null || command -v powershell &> /dev/null; then
    if (pwsh -Command "& {Write-Output 'PSH'}" &> /dev/null || powershell -Command "& {Write-Output 'PSH'}" &> /dev/null); then
        USER_SHELL = "powershell"
    fi
fi



#if on windows, need to verify if using wsl so 
#that a symbolic link can be created

shell_install() {
    case "$USER_SHELL" in
        *bash*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.bashrc
            echo ". ~/$SCRIPT_PATH" >> ~/.bashrc
            ;;
        *zsh*)
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.zshrc
            echo ". ~/$SCRIPT_PATH" >> ~/.zshrc
            ;;
        *powershell*)
            if [ -f "$PROFILE" ]; then
                echo "\$FT_EXE_PATH = \"$FT_EXE_PATH\"" >> "$PROFILE"
                echo ". '~/$SCRIPT_PATH'" >> "$PROFILE"
            fi
            ;;
        *)
            echo "Unrecognized shell. Please add the appropriate ftmain.sh to your shell's configuration file manually. Consider submitting a PR as well :)"
            exit 1
            ;;
    esac

    sudo cp $SCRIPT_PATH ~/

}



shell_install




