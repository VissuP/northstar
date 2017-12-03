local nsQL = require("nsQL")

function main()
local query = [[ INSERT INTO nssim.sampledata (rowid,id,mapdata) VALUES ( 'cbb33c38-e51b-4846-bd7c-4d6ba0c7b68d','cbb33c38-e51b-4846-bd7c-4d6ba0c7b68d',{'blah':'0x0a09092020097b0a090920200909226e616d65223a2022627269616e732d6275636b6574222c0a090920200909226372656174696f6e44617465223a2022323031372d30352d33315432313a30373a34342e3131315a220a09092020097d'});]]
	local source = {
		Protocol = "cassandra",
		Host = "cassandra1-dev2-northstar.mon-marathon-service.mesos,cassandra2-dev2-northstar.mon-marathon-service.mesos,cassandra3-dev2-northstar.mon-marathon-service.mesos,cassandra4-dev2-northstar.mon-marathon-service.mesos,cassandra5-dev2-northstar.mon-marathon-service.mesos",
		Port = "31814",
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
