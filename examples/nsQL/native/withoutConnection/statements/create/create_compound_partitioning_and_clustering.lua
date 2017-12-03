local nsQL = require("nsQL")

function main()
	local query = [[
		CREATE TABLE IF NOT EXISTS nssim.test4 (col1 text, col2 int, col3 double, col4 float, col5 text, PRIMARY KEY((col1, col2), col3, col4));
	]]
	local source = {
		Protocol = "cassandra",
		Host = "10.32.49.6",
		Port = "31838",
		Backend = "native"
	}
	processQuery(query, source, {})
end

function processQuery(query, source, options)
	local resp, err = nsQL.query(query, source, options)
	if(err ~= nil) then
		error(err)
	end
	return resp
end