#!/usr/bin/env bash


TARGET_DIR="$HOME/.local/share/fastTravelCLI"



# get the version
version_tag=$(git describe --tags --abbrev=0)

# get the logo
logo=$(bash -c "install/logo.sh")


exe_install() {
        
    sudo mkdir -p "$TARGET_DIR"
    echo "Created dir $TARGET_DIR"

    echo "Attempting to compile go binary..."
    go build -o fastTravel.exe -ldflags="-X 'github.com/osteensco/fastTravelCLI/ft.Version=$version_tag' -X 'github.com/osteensco/fastTravelCLI/ft.Logo=$logo'" 

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
            sed -i '/#FastTravel/,/#ftend/d' "$CONFIG"
            echo "#FastTravel" >> $CONFIG
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> $CONFIG
            echo ". $TARGET_DIR/ftmain.sh" >> $CONFIG
            echo "export FT_PATH=\"$TARGET_DIR/ftmain.sh\"" >> $CONFIG
            echo "#ftend" >> $CONFIG
            ;;
        *zsh*)
            SCRIPT_PATH="shells/bash/ftmain.sh"
            CONFIG=~/.zshrc
            sed -i '/#FastTravel/,/#ftend/d' "$CONFIG"
            echo "#FastTravel" >> $CONFIG
            echo "export FT_EXE_PATH=\"$FT_EXE_PATH\"" >> $CONFIG
            echo ". $TARGET_DIR/ftmain.sh" >> $CONFIG
            echo "export FT_PATH=\"$TARGET_DIR/ftmain.sh\"" >> $CONFIG
            echo "#ftend" >> $CONFIG
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
    
    echo "$logo"
    echo "$version_tag"
    echo "Config updated, please source $CONFIG or restart shell"
    
}



shell_install




