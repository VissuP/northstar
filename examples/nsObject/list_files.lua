local object = require("nsObject")
local output = require("nsOutput")

function main()
    local files, err = object.listFiles("test-bucket")
    if err ~= nil then
        error(err)
    end

    for i = 1, #files, 1 do
        output.print(files[i].key, " ", files[i].last_modified, " ", files[i].size, " ", files[i].etag, " ",
            files[i].storage_class, "\n")
    end
end