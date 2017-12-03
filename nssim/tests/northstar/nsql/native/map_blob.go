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

package native

import (
	"fmt"

	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.NSQLNativeMapBlob), NewNSQLNativeMapBlobTest)
}

//NSQLNativeMapBlobTest  implements the NSQL Native Map Blob test
type NSQLNativeMapBlobTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.NSQLFunctionality
}

func NewNSQLNativeMapBlobTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &NSQLNativeMapBlobTest{
		NorthstarApiBaseTest: nsapiBase,
		//Extend NorthstarApiBaseTest with NSQL functionality
		NSQLFunctionality: northstar.NewNSQLFunctionality(nsapiBase),
	}, nil
}

var (
	//steps

	//common variables
	mapBlobRowName = "NSQL Map Blob Test"
)

func (test *NSQLNativeMapBlobTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NSQL Native Map Blob Test")

	//Populate our sample data
	var generateTestDataStep utils.Step = "Generate test data."
	testData, err := test.populateTestData(logs)
	if err != nil {
		return logs.LogStep(generateTestDataStep, false, "Failed to populate test data with error: %s", err.Error())
	}

	tableName := "nssim.sampledata"
	wheres := "rowid = `" + fmt.Sprintf("%s", testData["rowid"]) + "`"
	logs.LogStep(generateTestDataStep, true, "Test data: %+v wheres: %s table: %s", testData, wheres, tableName)

	//Create row with our sample data
	var insertRowStep utils.Step = "Insert row"
	if err := test.InsertRow(logs, callbacks, tableName, testData); err != nil {
		return logs.LogStep(insertRowStep, false, "Error inserting row: %s", err.Error())
	}
	logs.LogStep(insertRowStep, true, "Success.")
	defer test.DeleteRow(logs, callbacks, tableName, wheres)

	//Retrieve row from cassandra
	var mapBlobFetchStep utils.Step = "Perform map blob fetch."
	mapFields := []northstar.NsqlField{
		{
			Column:   "mapdata",
			Field:    "blah",
			Subfield: "name",
			Alias:    "name",
		},
	}

	output, err := test.MapBlobFetch(logs, callbacks, tableName, wheres, mapFields)
	if err != nil {
		return logs.LogStep(mapBlobFetchStep, false, "Error getting row: %s", err.Error())
	}
	logs.LogStep(mapBlobFetchStep, true, "Success. Response: %+v", output)

	//Parse the row and make sure everything looks good
	var validateResultsStep utils.Step = "Validate results."
	table, err := test.TableToMap(logs, output.ExecutionResults, nil)
	if err != nil {
		return logs.LogStep(validateResultsStep, false, "Error, failed to parse results: %s", err.Error())
	}
	logs.LogDebug("Parse results: %+v", table)

	results := map[string]interface{}{
		"name": "brians-bucket",
	}
	err = test.VerifyResults(logs, results, table)
	if err != nil {
		return logs.LogStep(validateResultsStep, false, "Failed to verify inserted fields with error: %s", err.Error())
	}
	logs.LogStep(validateResultsStep, true, "Success.")

	return nil
}

//populateTestData populates the data used for the test
func (test *NSQLNativeMapBlobTest) populateTestData(logs *utils.Logger) (map[string]interface{}, error) {
	rowID, err := test.GetUUID()
	if err != nil {
		return nil, logs.LogError("Failed to get row ID with error: %s", err)
	}

	return map[string]interface{}{
		"rowid": rowID,
		"id":    rowID,
		"mapdata": map[string]interface{}{"blah": `
		  	{
		  		"name": "brians-bucket",
		  		"creationDate": "2017-05-31T21:07:44.111Z"
		  	}
		`},
		"name": mapBlobRowName,
	}, nil
}

func (test *NSQLNativeMapBlobTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	wheres := "name = '" + mapBlobRowName + "'"

	result, err := test.GetRow(logs, callbacks, tableName, wheres, northstar.NSQLOptions{AllowFiltering: true})
	if err != nil {
		return logs.LogError("Failed to get rows. %s", err.Error())
	}

	rows, err := test.TableToMap(logs, result.ExecutionResults, make(map[string]string))
	if err != nil {
		return logs.LogError("Could not parse rows. %s", err.Error())
	}

	for _, rowid := range rows["rowid"].Row {
		wheres = "rowid = '" + fmt.Sprintf("%s", rowid) + "'"
		if err := test.DeleteRow(logs, callbacks, tableName, wheres); err != nil {
			logs.LogError("Could not delete rows for wheres %s. %s", wheres, err.Error())
		}
	}
	return nil
}
