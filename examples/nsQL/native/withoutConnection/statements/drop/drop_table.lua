local nsQL = require("nsQL")

function main()
	local query = [[
		DROP TABLE nssim.test1;
	]]
	local source = {
		Protocol = "cassandra",
		Host = "10.32.49.6",
		Port = "31838",
		Backend = "native"
	}
	local options = {}
	processQuery(query, source, options)
end

function processQuery(query, source, options)
	local resp, err = nsQL.query(query, source, options)
	if(err ~= nil) then
		error(err)
	end
	return resp
end