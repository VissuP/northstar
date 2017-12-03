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
    local connection = createConnection(source)
    local result = processQuery(connection, query, options)
    teardownConnection(connection)
    return generateTable(result)
end

function createConnection(source)
    local connection, err = nsQL.connect(source)
    if(err ~= nil) then
        error(err)
    end
    return connection
end

function teardownConnection(connection)
    local err = connection:disconnect()
    if(err ~= nil) then
        error(err)
    end
end

function processQuery(connection, query, options)
    local resp, err = connection:query(query, options)
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