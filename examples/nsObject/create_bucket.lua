local object = require("nsObject")

function main()
    local err = object.createBucket("test-bucket")
    if err ~= nil then
        error(err)
    end
end