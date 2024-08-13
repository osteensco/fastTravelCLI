
case "$SHELL" in
    *bash*)
        CONFIG=~/.bashrc
        ;;
    *zsh*)
        CONFIG=~/.zshrc
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


if [ $# -ne 1 ]; then
    echo "shell found: $SHELL, profile: $CONFIG"
else
    echo "profile found: $CONFIG"
    CONFIG=$1
    echo "profile given: $CONFIG, attempting to cleanup $CONFIG..."
fi


if [ ! -f "$CONFIG" ]; then
    echo "File not found!"
    exit 1
fi

head -n -3 "$CONFIG" > temp_file && mv temp_file "$CONFIG"
    
