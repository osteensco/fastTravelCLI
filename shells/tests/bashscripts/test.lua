local maxSims

if #arg ~= 1 then
    maxSims = 15
else
    assert(tonumber(arg[1]), "Argument must be a number! Arg: " .. arg[1])
    maxSims = tonumber(arg[1])
end

local commands = { "]", "[", "..", "-", "hist" }

local generateSim = function(limit)
    local sim = {}
    local num = math.random(1, limit)
    for _ = 1, num do
        table.insert(sim, commands[math.random(1, #commands)])
    end
    return sim
end

local simulation = generateSim(maxSims)

for _, cmd in ipairs(simulation) do
    local handle = io.popen('bash -i -c "source ~/.bashrc; ft ' .. cmd .. '"')
    assert(handle, "handle cannot be nil")
    local output = handle:read('*a')
    handle:close()
    local formattedOutput = output:gsub('[\n\r]+', ' ')
    print(formattedOutput)
end
