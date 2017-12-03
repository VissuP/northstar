local object = require("nsObject")
local output = require("nsOutput")

function main()
    local buckets, err = object.listBuckets()
    if err ~= nil then
        error(err)
    end

    for i = 1, #buckets, 1 do
        output.print(buckets[i].name, " ", buckets[i].date, "\n")
    end
end