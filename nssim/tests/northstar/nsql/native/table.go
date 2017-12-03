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
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.NSQLNativeTable), NewNSQLNativeTableTest)
}

//NSQLNativeTableTest implements the Notebook NSQL Native Table test
type NSQLNativeTableTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.NSQLFunctionality
}

//NewNSQLNativeTableTest creates a new instance of the NSQLNativeTableTest
func NewNSQLNativeTableTest() (tests.Test, error) {
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
	createCompoundTableStep utils.Step = "Create Compound Table"
	deleteCompoundTableStep utils.Step = "Delete Compound Table"

	createCollectionTableStep utils.Step = "Create Clustered Table"
	deleteCollectionTableStep utils.Step = "Delete Clustered Table"

	createDirectiveTableStep utils.Step = "Create Directive Table"
	deleteDirectiveTableStep utils.Step = "Delete Directive Table"

	//common variables
	compoundTableName   = "nssim.compound." + config.Configuration.Environment
	collectionTableName = "nssim.collection." + config.Configuration.Environment
	directiveTableName  = "nssim.collection." + config.Configuration.Environment
)

//Execute exeutes the NSQL Native Crud Test
func (test *NSQLNativeTableTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NSQL Native Table Test")

	//compound partitioning and clustering
	output, err := test.NSQLFunctionality.ExecuteQuery(logs, callbacks, "CREATE TABLE IF NOT EXISTS `"+
		compoundTableName+"`CREATE TABLE IF NOT EXISTS nssim.test4 (col1 text, col2 int, col3 double, col4 float, col5 text, PRIMARY KEY((col1, col2), col3, col4));`", "")
	if err != nil {
		return logs.LogStep(createCompoundTableStep, false, "Failed to create table %s", compoundTableName)
	}
	logs.LogStep(createCompoundTableStep, true, "Successfully created table. %s. Output: %+v", tableName, output)

	output, err = test.NSQLFunctionality.DropTable(logs, callbacks, compoundTableName)
	if err != nil {
		logs.LogStep(deleteCompoundTableStep, false, "Failed to delete table %s. Output was: %+v", compoundTableName, output)
	}
	logs.LogStep(deleteCompoundTableStep, true, "Successfully deleted table %s.", compoundTableName)

	//clustering types
	output, err = test.NSQLFunctionality.ExecuteQuery(logs, callbacks, "CREATE TABLE IF NOT EXISTS `"+
		collectionTableName+"`CREATE TABLE IF NOT EXISTS nssim.test4 (col1 text, col2 int, col3 double, col4 float, col5 text, PRIMARY KEY((col1, col2), col3, col4));`", "")
	if err != nil {
		return logs.LogStep(createCollectionTableStep, false, "Failed to create table %s", collectionTableName)
	}
	logs.LogStep(createCompoundTableStep, true, "Successfully created table. %s. Output: %+v", tableName, output)

	output, err = test.NSQLFunctionality.DropTable(logs, callbacks, collectionTableName)
	if err != nil {
		logs.LogStep(deleteCollectionTableStep, false, "Failed to delete table %s. Output was: %+v", collectionTableName, output)
	}
	logs.LogStep(deleteCollectionTableStep, true, "Successfully deleted table %s.", collectionTableName)

	//directive
	output, err = test.NSQLFunctionality.ExecuteQuery(logs, callbacks, "CREATE TABLE IF NOT EXISTS `"+
		collectionTableName+"`CREATE TABLE IF NOT EXISTS nssim.test4 (col1 text, col2 int, col3 double, col4 float, col5 text, PRIMARY KEY((col1, col2), col3, col4));`", "")
	if err != nil {
		return logs.LogStep(createDirectiveTableStep, false, "Failed to create table %s", directiveTableName)
	}
	logs.LogStep(createDirectiveTableStep, true, "Successfully created table. %s. Output: %+v", tableName, output)

	output, err = test.NSQLFunctionality.DropTable(logs, callbacks, directiveTableName)
	if err != nil {
		logs.LogStep(deleteDirectiveTableStep, false, "Failed to delete table %s. Output was: %+v", directiveTableName, output)
	}
	logs.LogStep(deleteDirectiveTableStep, true, "Successfully deleted table %s.", directiveTableName)

	return nil
}

func (test *NSQLNativeTableTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")
	logs.LogInfo("No manual cleanup for this test.")
	return nil
}
