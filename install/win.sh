





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
            SCRIPT_PATH="shells/bash/ftmain.sh"
            CONFIG=~/.bashrc
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> $CONFIG
            echo ". ~/ftmain.sh" >> $CONFIG
            echo "export FT_PATH=\"~/ftmain.sh\"" >> $CONFIG
            ;;
        *zsh*)
            SCRIPT_PATH="shells/bash/ftmain.sh"
            CONFIG=~/.zshrc
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> $CONFIG
            echo ". ~/ftmain.sh" >> $CONFIG
            echo "export FT_PATH=\"~/ftmain.sh\"" >> $CONFIG
            ;;
        *)
            echo "
                Unrecognized shell. 
                Please add the appropriate ftmain.sh to your shell's configuration file manually. 
                Consider submitting a PR as well :)
            "
            exit 1
            ;;
    esac
    
    echo "Adding ftmain.sh to home folder"
    sudo cp $SCRIPT_PATH ~/.
    

    echo "fancy fastTravel ascii here"


    echo "Config updated, please source $CONFIG or restart shell"
    
}



shell_install




