


-- install_path specific to os
function test_install_script(install_path)
    local handle = io.popen(install_path .. " 2>&1")
    local result = handle:read("*a")
    local success = handle:close()
    return success, result
end

function test_ft_cmd(cmd)
    local handle = io.popen(cmd .. " 2>&1")
    local result = handle:read("*a")
    local success = handle:close()
    return success, result
end 

function cleanup(l)
    os.remove(os.getenv("FT_EXE_PATH"))
    os.remove(os.getenv("FT_PATH"))
    
    local profile_path = os.getenv("HOME") .. os.getenv("CONFIG")
    local temp = profile_path .. ".tmp"

    local lines = {}
    for line in io.lines(profile_path) do
        table.insert(lines, line)
    end

    if #lines >= 3 then 
        lines[#lines] = nil
        lines[#lines] = nil
        lines[#lines] = nil
    end

    --output to temp

end




