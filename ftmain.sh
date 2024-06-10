


ft() {


    output="$( "$FT_EXE_PATH" "$@" )"
    echo -E $output 
    if [ -e "$output" ]; then 
        cd "$output"
    fi

}



