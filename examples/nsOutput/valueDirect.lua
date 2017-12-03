local output = require("nsOutput")

function main()
    local value = {
        type = "int",
        value = "5"
    }

    local err = output.valueDirect(value)
    if err ~= nil then
        error(err)
    end
end