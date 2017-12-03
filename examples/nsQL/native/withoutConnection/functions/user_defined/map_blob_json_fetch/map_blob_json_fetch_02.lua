local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query =   [[
        SELECT  event.deviceid,
                MAP_BLOB_JSON_FETCH(event.fields, 'reportinginterval', 'value') as reportinginterval,
                MAP_BLOB_JSON_FETCH(event.fields, 'location', 'longitude') as longitude,
                MAP_BLOB_JSON_FETCH(event.fields, 'location', 'latitude') as latitude,
                MAP_BLOB_JSON_FETCH(event.fields, 'battery', '') as battery
        FROM    data.event
        WHERE   event.foreignid = 'e181197f-b9e6-6c4d-ed4b-6149584b9088'
        LIMIT   100;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.47.11.22",
        --Host = "cassandra1-dev2-ts.mon-marathon-service.mesos",
        Port = "9042",
        Backend = "native"
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