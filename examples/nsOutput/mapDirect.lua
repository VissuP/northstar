local output = require("nsOutput")

function main()
    local map = {
        type = "Map",
        center = {latitude = 1, longitude = 1},
        zoom = 8,
        items = {
            {
                label = "My Bike",
                locations = {
                    {latitude = 0.4, longitude = 1.6},
                    {latitude = 0.7, longitude = 1.3},
                    {latitude = 1, longitude = 1},
                    {latitude = 1.3, longitude = 0.7},
                    {latitude = 1.6, longitude = 0.4}}
            },
            {
                label = "My Car",
                locations = {
                    {latitude = 0.2, longitude = 1.6},
                    {latitude = 0.5, longitude = 1.3},
                    {latitude = 0.8, longitude = 1},
                    {latitude = 1.1, longitude = 0.7},
                    {latitude = 1.4, longitude = 0.4}}
            }
        }
    }

    local err = output.mapDirect(map)
    if err ~= nil then
        error(err)
    end
end