local output = require("nsOutput")

function main()
    for i=1,10001,1 do
        output.printf("%s", "a")
    end
end