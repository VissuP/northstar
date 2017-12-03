local nsQL = require("nsQL")

function main()
	local query = [[ INSERT INTO nssim.sampledata (id,json,rowid) VALUES ( '70bb549c-1c80-4596-a22f-fb926615aa5e','{"name": "brians-bucket","creationDate": "2017-05-31T21:07:44.111Z"}','70bb549c-1c80-4596-a22f-fb926615aa5e');]]
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