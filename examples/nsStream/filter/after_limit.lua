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

    err = stream:limit(5):filter(do_filter_things):start()
    if err ~= nil then
        error(err)
    end
end

function do_filter_things(data)
    if data[1] == 79 or data[1] == 111 then
        return true
    end
    return false
end