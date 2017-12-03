-- Filter stats.

local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query =   [[
        SELECT  event.foreignid, MAP_BLOB_JSON_FETCH(event.fields, 'reportinginterval', 'value') as onat
        FROM    data.event
        LIMIT   10;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.37.8.13",
        --Host = "cassandra1-dev2-ts.mon-marathon-service.mesos",
        Port = "9042",
        Backend = "spark"
    }
    local options = {}
    local result = processQuery(query, source, options)
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


