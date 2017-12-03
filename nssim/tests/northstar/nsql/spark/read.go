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
	tests.Register(tests.TestId(config.NSQLSparkRead), NewNSQLCrudTest)
}

//NSQLSparkReadTest implements the NSQL Read Operations test
type NSQLSparkReadTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.NSQLFunctionality
	spark *northstar.NSQLSparkFunctionality
}

//Types is used to massage data into a user's expected form. Database does not have a way to determine a column type so we do it manually.
var types map[string]string = map[string]string{
	"timevalue": "time",
	"datevalue": "date",
}

//NewNSQLCrudTest creates a new instance of the NSQLCrudTest
func NewNSQLCrudTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}
	return &NSQLSparkReadTest{
		NorthstarApiBaseTest: nsapiBase,
		//Extend NorthstarApiBaseTest with NSQL functionality
		NSQLFunctionality: northstar.NewNSQLFunctionality(nsapiBase),
		spark:             northstar.NewNSQLSparkFunctionality(nsapiBase),
	}, nil
}

var (
	//steps
	getRowStep utils.Step = "Get Row"

	//common variables
	sparkReadRowName = "Spark Read Test"
)

//Execute exeutes the NSQL Crud Test
func (test *NSQLSparkReadTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NS Spark Read Test")

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

	//Retrieve row from spark
	output, err := test.spark.GetRow(logs, callbacks, tableName, wheres, false)
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
	logs.LogStep(verifyResultsStep, true, "Success.")

	return nil
}

func (test *NSQLSparkReadTest) populateTestData(logs *utils.Logger) (map[string]interface{}, error) {
	rowID, err := test.GetUUID()
	if err != nil {
		return nil, logs.LogError("Failed to get row ID with error: %s", err)
	}

	return map[string]interface{}{
		"rowid":       rowID,
		"id":          rowID,
		"numvalue":    2017,
		"varintvalue": 662607,
		"maxvalue":    9223372036854775807,
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
		"name": sparkReadRowName,
	}, nil
}

//CleanupHangingResource cleans up any resources left over from previous test runs.
func (test *NSQLSparkReadTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	wheres := "name = '" + sparkReadRowName + "'"

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
