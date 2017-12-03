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

    err = stream:fold(do_fold_things, 0):fold(do_fold_things, 0):start()
    if err ~= nil then
        error(err)
    end
end

function do_fold_things(data, numLetters)
    return numLetters + #data
end