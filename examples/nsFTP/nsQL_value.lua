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
        SELECT  TCOUNT()
        FROM    inventory.utdevice
        WHERE   activatedon IS NOT NULL;
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

    local table = {
        columns = {"value"},
        rows = {{data.value}}
    }

    local data, err =  output.tableToCsv(table)
    if err ~= nil then
        error(err)
    end

    err = connection:store("/value.csv", data)
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