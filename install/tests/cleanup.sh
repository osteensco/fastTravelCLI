if [ $# -ne 1 ]; then
    echo "Usage: $0 <input_file>"
    exit 1
fi

input_file=$1

if [ ! -f "$input_file" ]; then
    echo "File not found!"
    exit 1
fi

head -n -3 "$input_file" > temp_file && mv temp_file "$input_file"
    
