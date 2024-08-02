


ft() {

    # start=`date +%s%N`
    
    temp_output=$(mktemp)
    
    "$FT_EXE_PATH" "$@" | tee "$temp_output"

    output="$(tail -n 1 "$temp_output")"

    if [ -d "$output" ]; then 
        cd "$output";
    fi

    rm "$temp_output"
    # end=`date +%s%N`
    # echo Execution time `expr $end - $start`
}



