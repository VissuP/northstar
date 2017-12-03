local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[
    				    SELECT  MAP_BLOB_JSON_FETCH(mapdata, 'blah','name') as name
    				    FROM    nssim.sampledata
    				    WHERE   rowid = `cbb33c38-e51b-4846-bd7c-4d6ba0c7b68d`;
    				]]
    				local source = {
    				    Protocol = "cassandra",
    				    Host = "10.47.3.10",
  						Port = "31814",
     				   Backend = "spark"
    				}
    				local result = processQuery(query, source, {})
    				return generateTable(result)
				end

				function processQuery(query, source, options)
   				 local resp, err = nsQL.query(query, source, options)
   				 if(err ~= nil) then
   				     error(err)
   				 end
   				 return resp
				end

				function generateTable(table)
 				   local out, err = nsOutput.table(table)
 				   if(err ~= nil) then
 				       error(err)
 				   end
  				  return out
				end
