local output = require("nsOutput")
local sftp = require("nsSFTP")

function main()
    local destination = {
        HostPort = "localhost:22",
        User = "test",
        Password = "wireless"
    }

    local connection, err = sftp.connect(destination)
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

    err = connection:store("/tmp/data.csv", data)
    if err ~= nil then
        error(err)
    end

    err = connection:disconnect()
    if err ~= nil then
        error(err)
    end
end