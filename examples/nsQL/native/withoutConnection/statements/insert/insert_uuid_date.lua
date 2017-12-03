local nsQL = require("nsQL")

function main()
    local query = [[
        INSERT INTO nssim.sampleData (rowId, id, datevalue)
        VALUES                       (
                                        'aca7ae94-1fc9-11e7-93ae-92361f001953',
                                        '9e3a6e50-1fc9-11e7-93ae-92361f001953',
                                        '2017-03-01'
                                     );
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.42.3.11",
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