local ftp = require("nsFTP")

function main()
    for i = 1, 4, 1 do
        local _, err = ftp.connect("10.37.2.6:21")
        if err ~= nil then
            error(err)
        end
    end
end