local nsQL = require("nsQL")

function main()
	local query = [[
		INSERT INTO nssim.sampledata	(rowid,id,mapdata)
		VALUES 							(
			'cbb33c38-e51b-4846-bd7c-4d6ba0c7b68d',
			'cbb33c38-e51b-4846-bd7c-4d6ba0c7b68d',
			{'blah':'0x0a09092020097b0a090920200909226e616d65223a2022627269616e732d6275636b6574222c0a090920200909226372656174696f6e44617465223a2022323031372d30352d33315432313a30373a34342e3131315a220a09092020097d'
			}
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