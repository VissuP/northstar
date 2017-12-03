local output = require("nsOutput")
local ftp = require("nsFTP")

function main()
    local connection, err = ftp.connect("10.37.2.6:21")
    if err ~= nil then
        error(err)
    end

    err = connection:login("test", "test")
    if err ~= nil then
        error(err)
    end

    local table = {
        columns = {"column1", "column2"},
        rows = {{1, 2}, {3, 4}, {5, 6}}
    }

    local data, err =  output.tableToCsv(table)
    if err ~= nil then
        error(err)
    end

    err = connection:store("/data.csv", data)
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