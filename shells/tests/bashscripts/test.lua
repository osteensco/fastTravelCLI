local commands = { "]", "[", "hist", "..", "-" }
for _, cmd in ipairs(commands) do
    -- may need to adjust if subprocess doesnt inherit sourced commands.sh and ftmain.sh
    -- this should work because we are adding this part to the .bashrc
    local handle = io.popen('bash -i -c "source ~/.bashrc; ft ' .. cmd .. '"')
    local output = handle:read('*a')
    handle:close()
    local formattedOutput = output:gsub('[\n\r]+', ' ')
    print(formattedOutput)
end
