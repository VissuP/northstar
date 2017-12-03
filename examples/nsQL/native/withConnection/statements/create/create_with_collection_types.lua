local nsQL = require("nsQL")

function main()
	local query = [[
		CREATE TABLE IF NOT EXISTS nssim.test6 (col1 text, col2 set<int>, col3 list<double>, col4 map<float, text>, PRIMARY KEY(col1));
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