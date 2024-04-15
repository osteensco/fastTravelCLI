


ft() {


    output="$( "$FT_EXE_PATH" "$@" )"
    echo $output 
    if [ -e "$output" ]; then 
        cd "$output"
    fi

}



