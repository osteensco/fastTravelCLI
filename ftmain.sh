


ft() {

    start=`date +%s%N`
    output="$( "$FT_EXE_PATH" "$@" )"
    echo -E $output; 
    if [ -e "$output" ]; then 
        cd "$output";
    fi
    end=`date +%s%N`
    echo Execution time `expr $end - $start`
}



