local nsQL = require("nsQL")

function main()
    local query = [[
        DELETE
        FROM    nssim.sampleData
        WHERE   rowId = 'aca7ae94-1fc9-11e7-93ae-92361f001953';
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.38.13.7",
        Port = "31838",
        Backend = "native"
    }
    processQuery(query, source, {})
end

function processQuery(query, source, options)
    local resp, err = nsQL.query(query, source, options)
    if(err ~= nil) then
        error(err)
    end
    return resp
end