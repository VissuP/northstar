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
	"fmt"
	"reflect"
	"strings"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/utils"
)

//NSQLSparkFunctionality provides support for NSQL capabilities
type NSQLSparkFunctionality struct {
	*NorthstarApiBaseTest
}

//NewNSQLSparkFunctionality creates a new instance of the NSQLSparkFunctionality library
func NewNSQLSparkFunctionality(nsAPIBase *NorthstarApiBaseTest) *NSQLSparkFunctionality {
	return &NSQLSparkFunctionality{
		NorthstarApiBaseTest: nsAPIBase,
	}
}

//typeToNSQL converts data so that can be used in an SQL statement
func (test *NSQLSparkFunctionality) typeToNSQL(value interface{}) (string, error) {
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

//VerifyResults verifies that the key/value pairs in expected values match the key/value pairs in actual values.
func (test *NSQLSparkFunctionality) VerifyResults(logs *utils.Logger, expectedValues map[string]interface{}, actualValues map[string][]interface{}) error {
	logs.LogInfo("verifyResults")

	for key, value := range expectedValues {
		actualValue, ok := actualValues[key]
		if !ok {
			return logs.LogError("Error, key: %s not set.", key)
		}

		logs.LogInfo("Checking value of %s. Value is: %v. Type is: %s.", key, actualValue[0], reflect.TypeOf(actualValue[0]))

		formattedValue, err := test.typeToNSQL(value)
		if err != nil {
			return logs.LogError("Error, cannot format value in NSQL format. Error was: %s", err.Error())
		}

		//Clean up quotes for strings
		formattedValue = strings.TrimPrefix(formattedValue, "'")
		formattedValue = strings.TrimSuffix(formattedValue, "'")

		if actualValues[key][0] != formattedValue {
			return logs.LogError("Error, key: %s value: %v does not match expected value: %v.", key, actualValues[key][0], formattedValue)
		}
	}

	return nil
}

//GetRow returns the row specified by the wheres in the specified table.
func (test *NSQLSparkFunctionality) GetRow(logs *utils.Logger,
	callbacks *utils.ThreadSafeMap,
	table string,
	wheres string,
	typed bool) (*northstarApiModel.Output, error) {

	logs.LogInfo("getRow")

	options := ""
	if typed {
		options = "ReturnTyped=true"
	}

	code := `local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[
    				    SELECT  *
    				    FROM    ` + table + `
    				    WHERE   ` + wheres + `;
    				]]
    				local source = {
    				    Protocol = "cassandra",
    				    Host = "` + config.Configuration.CassandraHost + `",
  						Port = "` + config.Configuration.CassandraPort + `",
     				    Backend = "spark"
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

	logs.LogInfo("Select row output: %+v", output)
	return output, nil
}

//JsonFetch retrieves the JSON fields specified
func (test *NSQLSparkFunctionality) JsonFetch(logs *utils.Logger,
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

	code := `local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[
    				    SELECT  ` + selects + `
    				    FROM    ` + table + `
    				    WHERE   ` + wheres + `;
    				]]
    				local source = {
    				    Protocol = "cassandra",
    				    Host = "` + config.Configuration.CassandraHost + `",
  						Port = "` + config.Configuration.CassandraPort + `",
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
				end`

	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return nil, err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return nil, err
	}

	logs.LogInfo("Select row output: %+v", output)
	return output, nil
}

//MapBlobFetch retrieves the Map fields specified
func (test *NSQLSparkFunctionality) MapBlobFetch(logs *utils.Logger,
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

	code := `local nsQL = require("nsQL")
				local nsOutput = require("nsOutput")

				function main()
    				local query = [[
    				    SELECT  ` + selects + `
    				    FROM    ` + table + `
    				    WHERE   ` + wheres + `;
    				]]
    				local source = {
    				    Protocol = "cassandra",
    				    Host = "` + config.Configuration.CassandraHost + `",
  						Port = "` + config.Configuration.CassandraPort + `",
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
				end`

	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return nil, err
	}

	output, err := validateExecutionResult(logs, *response)
	if err != nil {
		return nil, err
	}

	logs.LogInfo("Select row output: %+v", output)
	return output, nil
}
