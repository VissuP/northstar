local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[
    				    SELECT  JSON_FETCH(json, 'name') as name
    				    FROM    nssim.sampledata
    				    WHERE   rowid = `70bb549c-1c80-4596-a22f-fb926615aa5e`;
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
