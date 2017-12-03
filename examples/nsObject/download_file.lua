local object = require("nsObject")
local output = require("nsOutput")

function main()
    local data, err = object.downloadFile("test-bucket", "test-file-from-string")
    if err ~= nil then
        error(err)
    end

    output.print(data.ContentType)

    for i = 1, #data.Payload, 1 do
        output.print(data.Payload[i])
    end

    output.print("\n")

    data, err = object.downloadFile("test-bucket", "test-file-from-byte-array")
    if err ~= nil then
        error(err)
    end

    output.print(data.ContentType)

    for i = 1, #data.Payload, 1 do
        output.print(data.Payload[i])
    end
end