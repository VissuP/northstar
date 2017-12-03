local object = require("nsObject")

function main()
    local err = object.uploadFile("test-bucket", "test-file-from-string", "test-data", "text/plain")
    if err ~= nil then
        error(err)
    end
end