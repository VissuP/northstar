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
	tests.Register(tests.TestId(config.NSQLNativeCrud), NewNSQLNativeCrudTest)
}

//NSQLNativeCrudTest implements the Notebook NSQL Native Crud Operations test
type NSQLNativeCrudTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.NSQLFunctionality
}

//Types is used to massage data into a user's expected form. Database does not have a way to determine a column type so we do it manually.
var types map[string]string = map[string]string{
	"timevalue": "time",
	"datevalue": "date",
}

//NewNSQLNativeCrudTest creates a new instance of the NSQLNativeCrudTest
func NewNSQLNativeCrudTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}
	return &NSQLNativeCrudTest{
		NorthstarApiBaseTest: nsapiBase,
		//Extend NorthstarApiBaseTest with NSQL functionality
		NSQLFunctionality: northstar.NewNSQLFunctionality(nsapiBase),
	}, nil
}

var (
	//steps
	populateDataStep       utils.Step = "Generate test data"
	insertTestDataStep     utils.Step = "Insert test data"
	getRowStep             utils.Step = "Get row"
	updateRowStep          utils.Step = "Update row."
	getUpdatedRowStep      utils.Step = "Get updated row."
	validateUpdatedRowStep utils.Step = "Validate updated row."

	//common variables
	crudRowName = "NSQL CRUD Test"
	tableName   = "nssim.sampledata"
)

//Execute exeutes the NSQL Native Crud Test
func (test *NSQLNativeCrudTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NSQL Native CRUD Test")

	//Populate our sample data
	testData, err := test.populateTestData(logs)
	if err != nil {
		return logs.LogStep(populateDataStep, false, "Failed to generate test data with error: %s", err.Error())
	}

	wheres := "rowid = `" + fmt.Sprintf("%s", testData["rowid"]) + "` and id='" + fmt.Sprintf("%s", testData["id"]) + "'"
	logs.LogStep(populateDataStep, true, "Test data: %+v wheres: %s table: %s", testData, wheres, tableName)

	//Create row with our sample data
	if err := test.InsertRow(logs, callbacks, tableName, testData); err != nil {
		return logs.LogStep(insertTestDataStep, false, "Error inserting row: %s", err.Error())
	}
	defer test.DeleteRow(logs, callbacks, tableName, wheres)
	logs.LogStep(insertTestDataStep, true, "Success.")

	//Retrieve row from cassandra
	output, err := test.GetRow(logs, callbacks, tableName, wheres, northstar.NSQLOptions{})
	if err != nil {
		return logs.LogStep(getRowStep, false, "Error getting row: %s", err.Error())
	}
	logs.LogStep(getRowStep, true, "Response received: %+v", output)

	//Parse the row and make sure everything looks good
	var validateRowStep utils.Step = "Validate row."
	table, err := test.TableToMap(logs, output.ExecutionResults, types)
	if err != nil {
		return logs.LogStep(validateRowStep, false, "Error, failed to parse results: %s", err.Error())
	}

	logs.LogDebug("Parse results: %+v", table)
	err = test.VerifyResults(logs, testData, table)
	if err != nil {
		return logs.LogStep(validateRowStep, false, "Failed to verify inserted fields with error: %s", err.Error())
	}
	logs.LogStep(validateRowStep, true, "Success.")

	//Attempt to update a row
	updateFields := map[string]interface{}{
		"numvalue": "3000",
	}
	if err = test.UpdateRow(logs, callbacks, tableName, wheres, updateFields); err != nil {
		return logs.LogStep(updateRowStep, false, "Error updating row: %s", err.Error())
	}
	logs.LogStep(updateRowStep, true, "Success.")

	//Select the row again
	output, err = test.GetRow(logs, callbacks, tableName, wheres, northstar.NSQLOptions{})
	if err != nil {
		return logs.LogStep(getUpdatedRowStep, false, "Error getting row: %s", err.Error())
	}
	logs.LogStep(getUpdatedRowStep, true, "Success. Response: %+v", output)

	//Confirm our data got updated correctly
	table, err = test.TableToMap(logs, output.ExecutionResults, types)
	if err != nil {
		return logs.LogStep(validateUpdatedRowStep, false, "Error, failed to parse results: %s", err.Error())
	}
	err = test.VerifyResults(logs, updateFields, table)
	if err != nil {
		return logs.LogStep(validateUpdatedRowStep, false, "Failed to verify updated fields with error: %s", err.Error())
	}
	logs.LogStep(validateUpdatedRowStep, true, "Success.")

	return nil
}

func (test *NSQLNativeCrudTest) populateTestData(logs *utils.Logger) (map[string]interface{}, error) {
	rowID, err := test.GetUUID()
	if err != nil {
		return nil, logs.LogError("Failed to get row ID with error: %s", err)
	}

	return map[string]interface{}{
		"rowid":       rowID,
		"name":        crudRowName,
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
	}, nil
}

func (test *NSQLNativeCrudTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	wheres := "name = '" + crudRowName + "'"

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
