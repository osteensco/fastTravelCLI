
function ft
    
    set temp_output (mktemp)
        
    $FT_EXE_PATH $argv | tee $temp_output

    set output (tail -n 1 $temp_output)

    if test -d "$output"
        cd "$output"
    end

    rm $temp_output

end

