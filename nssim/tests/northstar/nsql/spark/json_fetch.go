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

package spark

import (
	"fmt"

	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.NSQLSparkJSONFetch), NewNSSparkJSONTest)
}

//NSQLSparkJSONTest  implements the NSQL Spark JSON fetch test
type NSQLSparkJSONTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.NSQLFunctionality
	spark *northstar.NSQLSparkFunctionality
}

var (
	//steps
	generateTestDataStep utils.Step = "Generate test data."
	insertTestDataStep   utils.Step = "Insert test data"
	jsonFetchStep        utils.Step = "JSON Fetch"
	verifyResultsStep    utils.Step = "Verify results"

	//common variables
	jsonFetchRowName = "NSQL Spark JSON Fetch Test"
	tableName        = "nssim.sampledata"
)

func NewNSSparkJSONTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &NSQLSparkJSONTest{
		NorthstarApiBaseTest: nsapiBase,
		//Extend NorthstarApiBaseTest with NSQL functionality
		NSQLFunctionality: northstar.NewNSQLFunctionality(nsapiBase),
		spark:             northstar.NewNSQLSparkFunctionality(nsapiBase),
	}, nil
}

func (test *NSQLSparkJSONTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NSQL Spark JSON Fetch Test")

	//Populate our sample data
	testData, err := test.populateTestData(logs)
	if err != nil {
		return logs.LogStep(generateTestDataStep, false, "Failed to populate test data with error: %s", err.Error())
	}

	tableName := "nssim.sampledata"
	wheres := "rowid = `" + fmt.Sprintf("%s", testData["rowid"]) + "`"
	logs.LogStep(generateTestDataStep, true, "Test data: %+v wheres: %s table: %s", testData, wheres, tableName)

	//Create row with our sample data
	if err := test.InsertRow(logs, callbacks, tableName, testData); err != nil {
		return logs.LogStep(insertTestDataStep, false, "Error inserting row: %s", err.Error())
	}
	logs.LogStep(insertTestDataStep, true, "Success.")
	defer test.DeleteRow(logs, callbacks, tableName, wheres)

	//Retrieve row from cassandra
	jsonFields := []northstar.NsqlField{
		{
			Column: "json",
			Field:  "name",
			Alias:  "name",
		},
	}
	output, err := test.spark.JsonFetch(logs, callbacks, tableName, wheres, jsonFields)
	if err != nil {
		return logs.LogStep(jsonFetchStep, false, "Error getting row: %s", err.Error())
	}
	logs.LogStep(jsonFetchStep, true, "Success. Response: %+v", output)

	//Parse the row and make sure everything looks good
	table, err := test.TableToMap(logs, output.ExecutionResults, nil)
	if err != nil {
		return logs.LogStep(verifyResultsStep, false, "Error, failed to parse results: %s", err.Error())
	}
	logs.LogDebug("Parse results: %+v", table)
	results := map[string]interface{}{
		"name": "brians-bucket",
	}
	err = test.VerifyResults(logs, results, table)
	if err != nil {
		return logs.LogStep(verifyResultsStep, false, "Failed to verify inserted fields with error: %s", err.Error())
	}
	logs.LogStep(verifyResultsStep, true, "Success")

	return nil
}

//populateTestData populates the data used for the test
func (test *NSQLSparkJSONTest) populateTestData(logs *utils.Logger) (map[string]interface{}, error) {
	rowID, err := test.GetUUID()
	if err != nil {
		return nil, logs.LogError("Failed to get row ID with error: %s", err)
	}

	return map[string]interface{}{
		"rowid": rowID,
		"id":    rowID,
		"json":  `{"name": "brians-bucket","creationDate": "2017-05-31T21:07:44.111Z"}`,
		"name":  jsonFetchRowName,
	}, nil
}

func (test *NSQLSparkJSONTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	wheres := "name = '" + jsonFetchRowName + "'"

	result, err := test.spark.GetRow(logs, callbacks, tableName, wheres, false)
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
