

T = {}

function T.set_env_vars(FT_EXE_PATH, FT_PATH, install_path)
    print("setting env variables...")
    T["FT_EXE_PATH"] = FT_EXE_PATH
    T["FT_PATH"] = FT_PATH
    T["FT_EXE_PATH_tmp"] = FT_EXE_PATH .. ".tmp"
    T["FT_PATH_tmp"] = FT_PATH .. ".tmp"
    T["install_path"] = install_path 
end

function T.prep()
    print("creating temp files...")
    os.rename(T["FT_EXE_PATH"],T["FT_EXE_PATH_tmp"])
    os.rename(T["FT_PATH"],T["FT_PATH_tmp"])
end

-- install_path is specific to OS
function T.test_install_script(install_path)
    print("testing install script...")
    print("sudo password: ")
    local handle = io.popen("bash \"" .. install_path .. "\" 2>&1")
    local result = handle:read("*a")
    assert(type(result) == "string", install_path .. " contents read is type " .. type(result) .. " expected string")
    local success = handle:close()
    return success, result
end

function T.cleanup()
    print("cleaning up...")
    os.remove(T["FT_EXE_PATH"])
    os.remove(T["FT_PATH"]) 
    os.rename(T["FT_EXE_PATH_tmp"], T["FT_EXE_PATH"])
    os.rename(T["FT_PATH_tmp"], T["FT_PATH"]) 

    local handle = io.popen("bash \"" .. "./install/tests/cleanup.sh" .. "\" 2>&1")
    local result = handle:read("*a")
    assert(type(result) == "string", "cleanup.sh contents read is type " .. type(result) .. " expected string")
    local success = handle:close()
    return success, result
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
    
    success, result = T.cleanup()
    if not success then
        print("cleanup script FAILED -> " .. result)
    else
        print("cleanup complete")
    end
    -- print("to cleanup rc file run cleanup script, usage: ./install/tests/cleanup.sh [profile]")
    
    print("install test completed")
end

return T
