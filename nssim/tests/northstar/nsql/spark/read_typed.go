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

var crudExpectedTypes map[string]string = map[string]string{
	"varintvalue": "int",
	"rowid":       "string",
	"json":        "string",
	"money":       "double",
	"timevalue":   "int",
	"name":        "string",
	"ip":          "string",
	"createdtime": "time",
	"floatvalue":  "double",
	"numvalue":    "int",
	"array":       "array[string]",
	"id":          "string",
	"maxvalue":    "int",
	"datevalue":   "time",
	"mapdata":     "map[string]blob",
	"data":        "blob",
}

func init() {
	tests.Register(tests.TestId(config.NSQLSparkTypedRead), NewNSQLSparkTypedReadTest)
}

//NotebookNSQLCrudTest implements the NSQL Spark Typed Read Operations test
type NSQLSparkTypedReadTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.NSQLFunctionality
	spark *northstar.NSQLSparkFunctionality
}

//NewNSQLSparkTypedReadTest creates a new instance of the NSQLSparkTypedReadTest
func NewNSQLSparkTypedReadTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}
	return &NSQLSparkTypedReadTest{
		NorthstarApiBaseTest: nsapiBase,
		//Extend NorthstarApiBaseTest with NSQL functionality
		NSQLFunctionality: northstar.NewNSQLFunctionality(nsapiBase),
		spark:             northstar.NewNSQLSparkFunctionality(nsapiBase),
	}, nil
}

var (
	//common variables
	sparkTypedReadRowName = "Spark Read Typed Test"
)

//Execute exeutes the NSQL Native Typed Crud Test
func (test *NSQLSparkTypedReadTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NSQL Native Typed CRUD Test")

	//Populate our sample data
	testData, err := test.populateTestData(logs)
	if err != nil {
		return logs.LogStep(generateTestDataStep, false, "Failed to populate test data with error: %s", err.Error())
	}

	tableName := "nssim.sampledata"
	wheres := "rowid = `" + fmt.Sprintf("%s", testData["rowid"]) + "` and id='" + fmt.Sprintf("%s", testData["id"]) + "'"
	logs.LogStep(generateTestDataStep, true, "Test data: %+v wheres: %s table: %s", testData, wheres, tableName)

	//Create row with our sample data
	if err := test.InsertRow(logs, callbacks, tableName, testData); err != nil {
		return logs.LogStep(insertTestDataStep, false, "Error inserting row: %s", err.Error())
	}
	logs.LogStep(insertTestDataStep, true, "Success.")
	defer test.DeleteRow(logs, callbacks, tableName, wheres)

	//Retrieve row from cassandra
	output, err := test.spark.GetRow(logs, callbacks, tableName, wheres, true)
	if err != nil {
		return logs.LogStep(getRowStep, false, "Error getting row: %s", err.Error())
	}
	logs.LogStep(getRowStep, true, "Success. Output: %+v. Results: %+v", output, output.ExecutionResults)

	//Parse the row and make sure everything looks good
	table, err := test.TableToMap(logs, output.ExecutionResults, types)
	if err != nil {
		return logs.LogStep(verifyResultsStep, false, "Error, failed to parse results: %s", err.Error())
	}
	logs.LogDebug("Parse results: %+v", table)
	err = test.VerifyResults(logs, testData, table)
	if err != nil {
		return logs.LogStep(verifyResultsStep, false, "Failed to verify inserted fields with error: %s", err.Error())
	}

	for column, rows := range table {
		expectedType, _ := crudExpectedTypes[column]
		if rows.DataType != expectedType {
			return logs.LogStep(verifyResultsStep, false, "Failed to verify types. Expected type %s not equal to actual type %s", expectedType, rows.DataType)
		}
	}
	logs.LogStep(verifyResultsStep, true, "Success.")

	return nil
}

func (test *NSQLSparkTypedReadTest) populateTestData(logs *utils.Logger) (map[string]interface{}, error) {
	rowID, err := test.GetUUID()
	if err != nil {
		return nil, logs.LogError("Failed to get row ID with error: %s", err)
	}

	return map[string]interface{}{
		"rowid":       rowID,
		"id":          rowID,
		"numvalue":    2017,
		"varintvalue": 662607,
		"maxvalue":    922337203685,
		"ip":          "127.0.0.1",
		"floatvalue":  3.1415,
		"money":       99.99,
		"timevalue":   "00:19:05",
		"datevalue":   "2049-03-04",
		"createdtime": "2017-06-05 12:00:05",
		"json": `
		[
		  		{
		  			"name": "brians-bucket",
		  			"creationDate": "2017-05-31T21:07:44.111Z"
		  		},
		  		{
		  			"name": "test-bucket",
		  			"creationDate": "2017-05-31T21:06:57.76Z"
		  		}
		  ]
		`,
		"name": sparkTypedReadRowName,
	}, nil
}

//CleanupHangingResource cleans up any resources left over from previous test runs.
func (test *NSQLSparkTypedReadTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	wheres := "name = '" + sparkTypedReadRowName + "'"

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
