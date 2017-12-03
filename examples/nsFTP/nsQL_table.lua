local output = require("nsOutput")
local ftp = require("nsFTP")
local nsQL = require("nsQL")

function main()
    local connection, err = ftp.connect("10.37.2.6:21")
    if err ~= nil then
        error(err)
    end

    err = connection:login("test", "test")
    if err ~= nil then
        error(err)
    end

    local query =   [[
        SELECT      COUNT(cellid) AS activity
        FROM        devicetxn.locationhistoryv2
        WHERE       msg_type = 'RetMsgStandard'
        GROUP BY    cellid;
    ]]
    local source = {
        Protocol = "cassandra",
        Host = "10.43.7.1",
        Port = "31055",
        Backend = "spark"
    }
    local options = {}
    local data, err = nsQL.query(query, source, options)
    if(err ~= nil) then
        error(err)
    end

    local data, err =  output.tableToCsv(data)
    if err ~= nil then
        error(err)
    end

    err = connection:store("/table.csv", data)
    if err ~= nil then
        error(err)
    end

    err = connection:logout()
    if err ~= nil then
        error(err)
    end

    err = connection:disconnect()
    if err ~= nil then
        error(err)
    end
end