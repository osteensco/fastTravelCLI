



SCRIPT_PATH="ftmain.sh"

TARGET_DIR=$("$HOME/AppData/Local"


exe_install() {
    
    ft_dir="$TARGET_DIR/fastTravel"

    sudo mkdir -p "$ft_dir"

    sudo mv ./fastTravel.exe "$ft_dir/" 

    echo "$ft_dir/fastTravel.exe"
}

FT_EXE_PATH=$(exe_install)

export PATH="$PATH:$FT_EXE_PATH"

# SHELL="echo $SHELL"


#if on windows, need to verify if using wsl so 
#that a symbolic link can be created

shell_install() {
    case "$SHELL" in
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
            echo "Unrecognized shell. Please add the script to your shell's configuration file manually."
            exit 1
            ;;
    esac

    sudo cp $SCRIPT_PATH ~/

}



shell_install




