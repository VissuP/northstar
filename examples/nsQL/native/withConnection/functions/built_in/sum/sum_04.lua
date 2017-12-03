local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query =   [[
        SELECT  SUM(invocations.elapsedtime) as sumtime
        FROM    account.invocations;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.32.49.6",
        Port = "31838",
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