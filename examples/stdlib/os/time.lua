local nsOutput = require("nsOutput")

function main()
    nsOutput.print(os.time(), ", ", os.time({year=1970, month=1, day=1, hour=0}))
end

