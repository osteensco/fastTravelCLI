





TARGET_DIR="/mnt/c/fastTravelCLI"



exe_install() {
    
    sudo mkdir -p "$TARGET_DIR"
    echo "Created dir $TARGET_DIR"

    echo "Attempting to compile go binary..."
    go build

    echo "Moving exe to $TARGET_DIR"
    sudo mv ./fastTravel.exe "$TARGET_DIR/fastTravel.exe" 

}


exe_install
FT_EXE_PATH="$TARGET_DIR/fastTravel.exe"



shell_install() {
    case "$SHELL" in
        *bash*)
            SCRIPT_PATH="./shells/bash/ftmain.sh"
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.bashrc
            echo ". ~/$SCRIPT_PATH" >> ~/.bashrc
            ;;
        *zsh*)
            SCRIPT_PATH="./shells/bash/ftmain.sh"
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> ~/.zshrc
            echo ". ~/$SCRIPT_PATH" >> ~/.zshrc
            ;;
        *)
            echo "Unrecognized shell. Please add the appropriate ftmain.sh to your shell's configuration file manually. Consider submitting a PR as well :)"
            exit 1
            ;;
    esac

    sudo cp $SCRIPT_PATH ~/

}



shell_install




