local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query =   [[
        SELECT  TCOUNT()
        FROM    account.invocations;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.38.13.7",
        Port = "31838",
        Backend = "native"
    }
    local options = {}
    local result = processQuery(query, source, options)
    return generateValue(result)
end

function processQuery(query, source, options)
    local resp, err = nsQL.query(query, source, options)
    if(err ~= nil) then
        error(err)
    end
    return resp
end

function generateValue(value)
    local out, err = nsOutput.value(value)
    if(err ~= nil) then
        error(err)
    end
    return out
end