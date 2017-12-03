local nsQL = require("nsQL")
local nsOutput = require("nsOutput")

function main()
    local query = [[
        SELECT  JSON_FETCH(stats, 'aof_enabled') as module,
                JSON_FETCH(stats, 'config_file') as etime,
                JSON_FETCH(stats, 'connected_slaves') as interval
        FROM    logging.stats
        WHERE   msgdate = '2017-09-07' AND `group` = 'northstar' AND appname = 'redis' AND interval = 13
        LIMIT   10;
    ]]

    local source = {
        Protocol = "cassandra",
        Host = "10.39.14.10",
        --Host = "cassandra1-log-dakota.mon-marathon-service.mesos",
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