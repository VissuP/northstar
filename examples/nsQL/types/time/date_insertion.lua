local nsQL = require("nsQL")

function main()
    local query = [[
        INSERT
        INTO    nssim.sampledata (rowid, id, datevalue)
        VALUES  ('aca7ae94-1fc9-11e7-93ae-92361f001955', '9e3a6e50-1fc9-11e7-93ae-92361f001955', '2049-03-04');
    ]]

    local source = {
        Protocol = "cassandra",
        Host = "10.6.15.24",
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

