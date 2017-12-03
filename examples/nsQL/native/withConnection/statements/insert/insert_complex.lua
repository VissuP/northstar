local nsQL = require("nsQL")

function main()
    local query = [[
        INSERT INTO nssim.sampleData    (
                                            rowId,
                                            id,
                                            createdtime,
                                            datevalue,
                                            timevalue,
                                            numvalue,
                                            maxvalue,
                                            varintvalue,
                                            name,
                                            floatvalue,
                                            money,
                                            ip,
                                            data,
                                            mapdata,
                                            array
                                        )
        VALUES                          (
                                            'aca7ae94-1fc9-11e7-93ae-92361f001953',
                                            '9e3a6e50-1fc9-11e7-93ae-92361f001953',
                                            '2017-06-05 00:00:00',
                                            '2017-06-05',
                                            '23:59:59',
                                            16,
                                            161616161616161616,
                                            1616,
                                            'somename',
                                            16.16,
                                            16.161616161616161616,
                                            '10.10.10.10',
                                            '0x01',
                                            {'blah':'0x00', 'bloh':'0x0F'},
                                            {'some_text_1', 'some_text_2'}
                                        );
    ]]

    local source = {
        Protocol = "cassandra",
        Host = "10.32.49.6",
        Port = "31838",
        Backend = "native"
    }
    local options = {}
    local connection = createConnection(source)
    processQuery(connection, query, options)
    teardownConnection(connection)
end

function createConnection(source)
    local connection, err = nsQL.connect(source)
    if(err ~= nil) then
        error(err)
    end
    return connection
end

function teardownConnection(connection)
    local err = connection:disconnect()
    if(err ~= nil) then
        error(err)
    end
end

function processQuery(connection, query, options)
    local resp, err = connection:query(query, options)
    if(err ~= nil) then
        error(err)
    end
    return resp
end