



# THIS WILL NEED REFACTORING FOR POWERSHELL ONCE INTEGRATED
TARGET_DIR="~/AppData/Local/fastTravelCLI"



exe_install() {
    
    sudo mkdir -p "$TARGET_DIR"
    echo "Created dir $TARGET_DIR"

    echo "Attempting to compile go binary..."
    go build -o fastTravel.exe

    echo "Moving exe to $TARGET_DIR"
    sudo mv ./fastTravel.exe "$TARGET_DIR/fastTravel.exe" 

}


exe_install
FT_EXE_PATH="$TARGET_DIR/fastTravel.exe"



# ONLY WINDOW SHELL TO BE SUPPORTED WILL BE POWERSHELL
shell_install() {
    case "$SHELL" in
        *bash*)
            SCRIPT_PATH="shells/bash/ftmain.sh"
            CONFIG=~/.bashrc
            echo "#FastTravel" >> $CONFIG
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> $CONFIG
            echo ". ~/ftmain.sh" >> $CONFIG
            echo "export FT_PATH=\"~/ftmain.sh\"" >> $CONFIG
            ;;
        *zsh*)
            SCRIPT_PATH="shells/bash/ftmain.sh"
            CONFIG=~/.zshrc
            echo "#FastTravel" >> $CONFIG
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
    
    echo "Adding ftmain.sh to $TARGET_DIR"
    sudo cp $SCRIPT_PATH $TARGET_DIR
    

    echo """
     __           _  _____                     _   ___   __   _____ 
    / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ 
   | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/ 
   |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  
   |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/  
                                                                    
    """


    echo "Config updated, please source $CONFIG or restart shell"
    
}



shell_install




