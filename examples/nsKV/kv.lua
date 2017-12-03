local nsKV = require("nsKV")
local output = require("nsOutput")

function main()
    err = nsKV.set("key1", 0, 10000)
    if err ~= nil then
        error(err)
    end

    val, err = nsKV.get("key1")
    if err ~= nil then
        error(err)
    end
    output.printf("Value retrieved: %v", val)

    val, err = nsKV.incr("key1")
    if err ~= nil then
        error(err)
    end
    output.printf("Value incremented: %v", val)

    val, err = nsKV.del("key1")
    if err ~= nil then
        error(err)
    end
end
