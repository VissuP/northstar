local object = require("nsObject")

function main()
    local err = object.deleteFile("test-bucket", "test-file-from-byte-array")
    if err ~= nil then
        error(err)
    end
end