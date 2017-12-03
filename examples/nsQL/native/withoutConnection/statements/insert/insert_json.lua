local nsQL = require("nsQL")

function main()
	local query = [[ INSERT INTO nssim.sampledata (id,json,rowid) VALUES ( '70bb549c-1c80-4596-a22f-fb926615aa5e','{"name": "brians-bucket","creationDate": "2017-05-31T21:07:44.111Z"}','70bb549c-1c80-4596-a22f-fb926615aa5e');]]
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
