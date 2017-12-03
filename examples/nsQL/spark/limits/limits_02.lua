local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query = [[
        SELECT  *
        FROM    account.invocations
        LIMIT   5;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.38.13.7",
        Port = "31838",
        Backend = "spark"
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

