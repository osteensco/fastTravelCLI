

T = {}

function T.set_env_vars(FT_EXE_PATH, FT_PATH, install_path)
    T["FT_EXE_PATH"] = FT_EXE_PATH
    T["FT_PATH"] = FT_PATH
    T["FT_EXE_PATH_tmp"] = FT_EXE_PATH .. ".tmp"
    T["FT_PATH_tmp"] = FT_PATH .. ".tmp"
    T["install_path"] = install_path 
end

function T.prep()
    os.rename(T["FT_EXE_PATH"],T["FT_EXE_PATH_tmp"])
    os.rename(T["FT_PATH"],T["FT_PATH_tmp"])
end

-- install_path is specific to OS
function T.test_install_script(install_path)
    local handle = io.popen("bash \"" .. install_path .. "\" 2>&1")
    local result = handle:read("*a")
    assert(type(result) == "string", install_path .. " contents read is type " .. type(result) .. " expected string")
    local success = handle:close()
    return success, result
end

function T.cleanup()
    os.remove(T["FT_EXE_PATH"])
    os.remove(T["FT_PATH"]) 
    os.rename(T["FT_EXE_PATH_tmp"], T["FT_EXE_PATH"])
    os.rename(T["FT_PATH_tmp"], T["FT_PATH"]) 
end

function T.main()
    local success, result
   
    T.prep()

    success, result = T.test_install_script(T["install_path"])
    if not success then
        print("install script - FAIL - " .. result)
    else
        print("install script - success")
    end
    
    T.cleanup()
    print("to cleanup rc file run cleanup script, usage: ./install/tests/cleanup.sh [profile]")
end

return T
