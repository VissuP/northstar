local output = require("nsOutput")

function main()
    local value = {
        type = "int",
        value = "5"
    }

    local out, err = output.value(value)
    if err ~= nil then
        error(err)
    end

    return out
end