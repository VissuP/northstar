local nsOutput = require("nsOutput")

function main()
    nsOutput.print(os.difftime(os.time(), os.time()))
end

