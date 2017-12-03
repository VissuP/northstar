local output = require("nsOutput")

function main()
    local table = {
        columns = {"column1", "column2"},
        rows = {{1, 2}, {3, 4}, {5, 6}}
    }

    local err = output.tableDirect(table)
    if err ~= nil then
        error(err)
    end
end