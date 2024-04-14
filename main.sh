

# ft() {
#     while IFS= read -r line; do
#         output="$line"
#     done < <("$FT_EXE_PATH" "$@")
                            
#     echo "$output" 
                                    
#     if [ -e "$output" ]; then 
#         cd "$output"
#     fi
# }



ft() {


    output="$( "$FT_EXE_PATH" "$@" )"
    echo $output 
    if [ -e "$output" ]; then 
        cd "$output"
    fi

}


