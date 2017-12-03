-- IMPORTANT: The correct flow of this example results in an error.

local nsStream = require("nsStream")

local kafkaSource = {
    Topic = "nsstream-test",
    Brokers = "10.37.13.6:9092",
    ZK = "10.44.6.3:2181"
}

function main()
    local stream, err = nsStream.create("kafka", "description", kafkaSource)
    if err ~= nil then
        error(err)
    end

    err = stream:fold(do_fold_things, 0):foreach(do_foreach_things):start()
    if err ~= nil then
        error(err)
    end
end

function do_fold_things(data, numLetters)
    return numLetters + #data
end

function do_foreach_things(data)
    local nsOutput = require("nsOutput")
    local asString = ""
    for i=1, #data, 1 do
        asString = asString .. string.char(data[i])
    end
    nsOutput.print(asString)
end