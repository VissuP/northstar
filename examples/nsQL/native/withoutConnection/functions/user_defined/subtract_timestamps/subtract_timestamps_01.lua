local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query = [[
        SELECT  SUBTRACT_TIMESTAMPS(finishedon, createdon)
        FROM    account.invocations;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.38.13.7",
        Port = "31838",
        Backend = "native"
    }

    local result = processQuery(query, source, {})
    return generateTable(result)
end

function processQuery(query, source, options)
    local resp, err = nsQL.query(query, source, options)
    if(err ~= nil) then
        error(err)
    end
    return resp
end

function generateTable(table)
    local out, err = nsOutput.table(table)
    if(err ~= nil) then
        error(err)
    end
    return out
end

