local object = require("nsObject")

function main()
    local err = object.uploadFile("test-bucket", "test-file-from-byte-array", {106, 107, 108}, "text/plain")
    if err ~= nil then
        error(err)
    end
end