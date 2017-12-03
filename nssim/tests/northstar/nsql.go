/*
Copyright (C) 2017 Verizon. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package northstar

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/utils"
)

//NSQLFunctionality provides support for NSQL capabilities
type NSQLFunctionality struct {
	*NorthstarApiBaseTest
}

//NewNSQLFunctionality creates a new instance of the NSQLFunctionality library
func NewNSQLFunctionality(nsAPIBase *NorthstarApiBaseTest) *NSQLFunctionality {
	return &NSQLFunctionality{
		NorthstarApiBaseTest: nsAPIBase,
	}
}

//typeToNSQL converts data so that can be used in an SQL statement
func (test *NSQLFunctionality) typeToNSQL(value interface{}) (string, error) {
	//Place strings in quotes. Otherwise include type as is
	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return `'` + fmt.Sprintf("%s", value) + `'`, nil
	case reflect.Slice:
		values, ok := value.([]interface{})
		if !ok {
			return "", fmt.Errorf("Error, could not parse map to []interface{}.")
		}
		data := "{"
		for _, entry := range values {
			data += "'" + fmt.Sprintf("%v", entry) + "',"
		}
		data = strings.TrimSuffix(data, ",")
		data += "}"
		return data, nil
	case reflect.Map:
		values, ok := value.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("Error, could not parse map to map[string]interface{}.")
		}
		data := "{"
		for key, entry := range values {
			data += "'" + key + "':" + fmt.Sprintf("'%v'", entry) + ","
		}
		data = strings.TrimSuffix(data, ",")
		data += "}"

		return data, nil
	default:
		return fmt.Sprintf("%v", value), nil
	}

	return "", fmt.Errorf("Error, escaped type switch statement in typeToNSQL")
}

//InsertRow will insert the contents of values (map of column->value into the specified table.)
func (test *NSQLFunctionality) InsertRow(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string,
	values map[string]interface{}) error {
	logs.LogInfo("InsertRow")

	var columns string
	var data string
	for column, value := range values {
		columns += column + ","

		entry, err := test.typeToNSQL(value)
		if err != nil {
			return logs.LogError("Error, cannot parse value: %v to NSQL statement: %s", value, err.Error())
		}
		data += entry + `,`
	}
	columns = strings.TrimSuffix(columns, ",")
	data = strings.TrimSuffix(data, ",")

	code := `local nsQL = require("nsQL")

					function main()
   				 	local query = [[ INSERT INTO ` + table + ` (` + columns + `) VALUES ( ` + data + `);]]
       				 	local source = {
       				 	    Protocol = "cassandra",
        				 	   Host = "` + config.Configuration.CassandraHost + `",
         				 	  Port = "` + config.Configuration.CassandraPort + `",
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
   				 	end`
	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return err
	}

	logs.LogInfo("Insert row output: %+v", output)
	return nil
}

//DeleteRow will delete the row specified via wheres from the specfied table.
func (test *NSQLFunctionality) DeleteRow(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string, wheres string) error {
	logs.LogInfo("DeleteRow")

	code := `local nsQL = require("nsQL")

				function main()
    				local query = [[
 				       DELETE
 				        FROM    ` + table + `
 				       WHERE   ` + wheres + `;
 				   ]]
 				   local source = {
 				      Protocol = "cassandra",
        			  Host = "` + config.Configuration.CassandraHost + `",
  				      Port = "` + config.Configuration.CassandraPort + `",
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
				end`

	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return err
	}

	logs.LogInfo("Delete row output: %+v", output)
	return nil
}

type NSQLOptions struct {
	Typed          bool
	AllowFiltering bool
}

func (test *NSQLFunctionality) ExecuteQuery(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap, query string, options string) (*northstarApiModel.Output, error) {
	logs.LogInfo("ExecuteQuery")

	code := `local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[` + query + `]]
    				local source = {
    				   Protocol = "cassandra",
    				   Host = "` + config.Configuration.CassandraHost + `",
  				       Port = "` + config.Configuration.CassandraPort + `",
     				   Backend = "native"
    				}
    				local result = processQuery(query, source, { ` + options + `})
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
				end`

	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return nil, err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//DeleteTable
func (test *NSQLFunctionality) DropTable(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string) (*northstarApiModel.Output, error) {
	logs.LogInfo("DropTable")

	code := `local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[ DROP TABLE ` + table + `]]
    				local source = {
    				   Protocol = "cassandra",
    				   Host = "` + config.Configuration.CassandraHost + `",
  				       Port = "` + config.Configuration.CassandraPort + `",
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
				end`

	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return nil, err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return nil, err
	}

	return output, nil
}

//GetRow returns the row specified by the wheres in the specified table.
func (test *NSQLFunctionality) GetRow(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string,
	wheres string,
	nsqlOptions NSQLOptions) (*northstarApiModel.Output, error) {

	logs.LogInfo("GetRow")

	options := ""
	if nsqlOptions.Typed {
		options += "ReturnTyped=true"
	}

	if nsqlOptions.AllowFiltering {
		if len(options) > 0 {
			options += ", "
		}
		options += "AllowFiltering=true"
	}

	query := `SELECT  *
    		FROM    ` + table + `
    		WHERE   ` + wheres + `;`

	output, err := test.ExecuteQuery(logs, callbacks, query, options)
	if err != nil {
		return nil, err
	}

	logs.LogInfo("Select row output: %+v", output)
	return output, nil
}

//UpdateRow updates the  the column:value pairs in the row specified via the wheres in the specified table
func (test *NSQLFunctionality) UpdateRow(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string,
	wheres string,
	values map[string]interface{}) error {
	logs.LogInfo("updateRow")

	var data string
	for column, value := range values {

		entry, err := test.typeToNSQL(value)
		if err != nil {
			return logs.LogError("Error, cannot parse value: %v to NSQL statement: %s", value, err.Error())
		}
		data = column + `=` + fmt.Sprintf("%v", entry)
	}
	data = strings.TrimSuffix(data, ",")

	code := `local nsQL = require("nsQL")

			function main()
			    local query = [[
  			      UPDATE  ` + table + `
  			      SET     ` + data + `
  			      WHERE   ` + wheres + `;
  			  ]]
  			  local source = {
  			      Protocol = "cassandra",
  			      Host = "` + config.Configuration.CassandraHost + `",
  			      Port = "` + config.Configuration.CassandraPort + `",
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
			end`
	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return err
	}

	logs.LogInfo("Update row output: %+v", output)
	return nil
}

//validateExecutionResult validates the execution result
func validateExecutionResult(logs *utils.Logger, response json.RawMessage) (*northstarApiModel.Output, error) {
	output := northstarApiModel.Output{}

	logs.LogInfo("Validating response: %s", response)
	err := json.Unmarshal([]byte(response), &output)
	if err != nil {
		return nil, logs.LogError("Error, failed to unmarshal output with error: %s", err)
	}

	logs.LogInfo("Execution output: %+v", output)
	if output.Status != northstarApiModel.OutputSuccessStatus {
		return &output, logs.LogError("Failed to execute with error: %s.", output.StatusDescription)
	}

	return &output, nil
}

//VerifyResults verifies that the key/value pairs in expected values match the key/value pairs in actual values.
func (test *NSQLFunctionality) VerifyResults(logs *utils.Logger, expectedValues map[string]interface{}, actualValues map[string]tableRow) error {
	logs.LogInfo("verifyResults")

	for key, value := range expectedValues {
		actualValue, ok := actualValues[key]
		if !ok {
			return logs.LogError("Error, key: %s not set.", key)
		}

		logs.LogInfo("Checking value of %s. Value is: %v. Type is: %s.", key, actualValue.Row[0], reflect.TypeOf(actualValue.Row[0]))

		formattedValue, err := test.typeToNSQL(value)
		if err != nil {
			return logs.LogError("Error, cannot format value in NSQL format. Error was: %s", err.Error())
		}

		//Clean up quotes for strings
		formattedValue = strings.TrimPrefix(formattedValue, "'")
		formattedValue = strings.TrimSuffix(formattedValue, "'")

		if actualValues[key].Row[0] != formattedValue {
			return logs.LogError("Error, key: %s value: %v does not match expected value: %v.", key, actualValues[key].Row[0], formattedValue)
		}
	}

	return nil
}

//NsqlField defines the JSON field to select
type NsqlField struct {
	Column   string
	Field    string
	Subfield string
	Alias    string
}

//JsonFetch retrieves the JSON fields specified
func (test *NSQLFunctionality) JsonFetch(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string, wheres string,
	fields []NsqlField) (*northstarApiModel.Output, error) {
	logs.LogInfo("JSONFetch")

	var selects string
	for _, field := range fields {
		selects = `JSON_FETCH(` + field.Column + `, '` + field.Field + `')`
		if field.Alias != "" {
			selects += " as " + field.Alias
		}
		selects += ","
	}
	selects = strings.TrimSuffix(selects, ",")

	query := `SELECT  ` + selects + `
	FROM    ` + table + `
	WHERE   ` + wheres + `;`

	output, err := test.ExecuteQuery(logs, callbacks, query, "")
	if err != nil {
		return nil, err
	}
	logs.LogInfo("Select row output: %+v", output)
	return output, nil
}

//MapBlobFetch retrieves the Map fields specified
func (test *NSQLFunctionality) MapBlobFetch(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string,
	wheres string,
	fields []NsqlField) (*northstarApiModel.Output, error) {
	logs.LogInfo("MapBlobFetch")

	var selects string
	for _, field := range fields {
		selects = `MAP_BLOB_JSON_FETCH(` + field.Column + `, '` + field.Field + `','` + field.Subfield + `')`
		if field.Alias != "" {
			selects += " as " + field.Alias
		}
		selects += ","
	}
	selects = strings.TrimSuffix(selects, ",")

	query := `SELECT  ` + selects + `
	FROM    ` + table + `
	WHERE   ` + wheres + `;`

	output, err := test.ExecuteQuery(logs, callbacks, query, "")
	if err != nil {
		return nil, err
	}
	logs.LogInfo("Select row output: %+v", output)
	return output, nil

	logs.LogInfo("Select row output: %+v", output)
	return output, nil
}
