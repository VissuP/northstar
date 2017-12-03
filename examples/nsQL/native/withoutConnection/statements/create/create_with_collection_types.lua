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
	processQuery(query, source, {})
end

function processQuery(query, source, options)
	local resp, err = nsQL.query(query, source, options)
	if(err ~= nil) then
		error(err)
	end
	return resp
end