local output = require("nsOutput")

function main()
    local table = {}
    table["test"] = "hello world"

    local out, err = output.html("<html><b>" .. table["test"] .. "</b></html>")
    if err ~= nil then
        error(err)
    end

    return out
end