local output = require("nsOutput")

function main()
    local table = {}
    table["test"] = "hello world"

    local err = output.htmlDirect("<html><b>" .. table["test"] .. "</b></html>")
    if err ~= nil then
        error(err)
    end
end