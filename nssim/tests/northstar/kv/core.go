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

package nsql

import (
	"fmt"
	"strings"
	"time"

	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.NSKVCore), NewKVCoreTest)
}

//NotebookKVCoreTest implements the Notebook Key Value Core test
type NotebookKVCoreTest struct {
	*northstar.NorthstarApiBaseTest
	*northstar.KVFunctionality
}

//NewNSQLCrudTest creates a new instance of the KVCoreTest
func NewKVCoreTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}
	return &NotebookKVCoreTest{
		NorthstarApiBaseTest: nsapiBase,
		//Extend NorthstarApiBaseTest with key value functionality
		KVFunctionality: northstar.NewKVFunctionality(nsapiBase),
	}, nil
}

var (
	//Steps in this test
	setTestKeyStep                 utils.Step = "Setting Test Key Value"
	retrieveValueStep              utils.Step = "Retrieving Value"
	setValueNotExistsButExistsStep utils.Step = "Setting Value If Not Exists Does Exist"
	SetValueNotExistsDNEStep       utils.Step = "Setting Value If Not Exists Value Does Not Exist"
	incrementValueStep             utils.Step = "Incrementing Value"
	decrementValueStep             utils.Step = "Decrementing Value"
	keyExistsStep                  utils.Step = "Checking If Key Exists Key Exists"
	KeyExistsDNEStep               utils.Step = "Checking If Key Exists Key Does Not Exist"
	ExpireKeyStep                  utils.Step = "ExpireKey"
	CheckKeyExistsStep             utils.Step = "Checking If Key Exists Key Should No Longer Exist"
	AppendStep                     utils.Step = "Appending To Value"
	SetRangeStep                   utils.Step = "Set Range"
	GetRangeStep                   utils.Step = "Get Range"
	SetHashStep                    utils.Step = "Setting Hash"
	GetHashStep                    utils.Step = "Getting Hash"
	GetAllHashStep                 utils.Step = "Getting All Hash Values"

	//Common variables across several tests
	kvSetKey   = "nssimCoreTest"
	kvNXSetKey = "nssimCoreTest2"
	hashKey    = "nssimTestHash"
	hashField  = "hash"
	hashValue  = "testHash"
)

//Execute executes the NS Key Value Core Test
func (test *NotebookKVCoreTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- NS KV Core Test")

	kvSetValue := 0
	err := test.KVSet(logs, callbacks, kvSetKey, fmt.Sprintf("%v", kvSetValue), 100)
	if err != nil {
		return logs.LogStep(setTestKeyStep, false, "Error, could not set value for key %s: %s", kvSetKey, err.Error())
	}

	logs.LogStep(setTestKeyStep, true, "Success")
	defer test.KVDel(logs, callbacks, kvSetKey)

	returnedValue, err := test.KVGet(logs, callbacks, kvSetKey)
	if err != nil {
		return logs.LogStep(retrieveValueStep, false, "Error, could not get value for key %s: %s", kvSetKey, err.Error())
	}

	if returnedValue != fmt.Sprintf("%v", kvSetValue) {
		return logs.LogStep(retrieveValueStep, false, "Error, returned value from KV Store  (%v) does not equal original value (%v).", returnedValue, kvSetValue)
	}

	logs.LogStep(retrieveValueStep, true, "Success")

	success, err := test.KVSetNx(logs, callbacks, kvSetKey, "nssimTesting456", 100)
	if err != nil {
		return logs.LogStep(setValueNotExistsButExistsStep, false, "Error, KVSetNx (key exists) failed with error: %s", err.Error())
	}
	if success {
		return logs.LogStep(setValueNotExistsButExistsStep, false, "Error, value set when key exists")
	}
	logs.LogStep(setValueNotExistsButExistsStep, true, "Success")

	success, err = test.KVSetNx(logs, callbacks, kvNXSetKey, "nssimTesting456", 100)
	if err != nil {
		return logs.LogStep(SetValueNotExistsDNEStep, false, "Error, KVSetNx (key does not exist) failed with error: %s", err.Error())
	}
	if !success {
		return logs.LogStep(SetValueNotExistsDNEStep, false, "Error, failed to set value when key does not exist.")
	}
	logs.LogStep(SetValueNotExistsDNEStep, true, "Success")
	defer test.KVDel(logs, callbacks, kvNXSetKey)

	value, err := test.KVIncr(logs, callbacks, kvSetKey)
	if err != nil {
		return logs.LogStep(incrementValueStep, false, "Error, failed to increment key:%s with error: %s", kvNXSetKey, err.Error())
	}

	kvSetValue++
	if value != kvSetValue {
		return logs.LogStep(incrementValueStep, false, "Error, expected increment value %d not equal to %d", value, kvSetValue)
	}

	logs.LogStep(incrementValueStep, true, "Success")

	value, err = test.KVDecr(logs, callbacks, kvSetKey)
	if err != nil {
		return logs.LogStep(decrementValueStep, false, "Error, failed to increment key:%s with error: %s", kvNXSetKey, err.Error())
	}

	kvSetValue--
	if value != kvSetValue {
		return logs.LogStep(decrementValueStep, false, "Error, expected increment value %d not equal to %d", value, kvSetValue)
	}
	logs.LogStep(decrementValueStep, true, "Success")

	exists, err := test.KVExists(logs, callbacks, kvSetKey)
	if err != nil {
		return logs.LogStep(keyExistsStep, false, "Error checking if key %s exists: %s", kvSetKey, err.Error())
	}
	if !exists {
		return logs.LogStep(keyExistsStep, false, "Error key %s does not exist.", kvSetKey)
	}
	logs.LogStep(keyExistsStep, true, "Success")

	exists, err = test.KVExists(logs, callbacks, "doesNotExist")
	if err != nil {
		return logs.LogStep(KeyExistsDNEStep, false, "Error checking if key %s exists: %s", kvSetKey, err.Error())
	}
	if exists {
		return logs.LogStep(KeyExistsDNEStep, false, "Error key %s exists when it should not.", kvSetKey)
	}
	logs.LogStep(KeyExistsDNEStep, true, "Success")

	set, err := test.KVExpire(logs, callbacks, kvSetKey, 3)
	if err != nil {
		return logs.LogStep(ExpireKeyStep, false, "Error, could not set expire on key: %s. Error was: %s.", kvSetKey, err.Error())
	}
	if !set {
		return logs.LogStep(ExpireKeyStep, false, "Error, expire was not set on key: %s", kvSetKey)
	}
	logs.LogStep(ExpireKeyStep, true, "Success")

	time.Sleep(4 * time.Second)

	exists, err = test.KVExists(logs, callbacks, kvSetKey)
	if err != nil {
		return logs.LogStep(CheckKeyExistsStep, false, "Error checking if key %s exists: %s", kvSetKey, err.Error())
	}
	if exists {
		return logs.LogStep(CheckKeyExistsStep, false, "Error key %s exists when it should be expired.", kvSetKey)
	}
	logs.LogStep(CheckKeyExistsStep, true, "Success")

	length, err := test.KVAppend(logs, callbacks, kvSetKey, "appended")
	if err != nil {
		return logs.LogStep(AppendStep, false, "Error, failed to append to key %s with error: %s", kvSetKey, err.Error())
	}

	appendedValue, err := test.KVGet(logs, callbacks, kvSetKey)
	if err != nil {
		return logs.LogStep(AppendStep, false, "Error, failed to get new value of key: %s with error: %s", kvSetKey, err.Error())
	}
	logs.LogStep(AppendStep, true, "New value after append: %s", appendedValue)

	rangeValue := "test"
	length, err = test.KVSetRange(logs, callbacks, kvSetKey, 0, rangeValue)
	if err != nil {
		return logs.LogStep(SetRangeStep, false, "Error, failed to set range for key: %s with error: %s", kvSetKey, err.Error())
	}
	logs.LogStep(SetRangeStep, true, "New length: %d", length)

	rangeResult, err := test.KVGetRange(logs, callbacks, kvSetKey, 0, len(rangeValue)-1)
	if err != nil {
		return logs.LogStep(GetRangeStep, false, "Error, failed to get range for key: %s with error: %s", kvSetKey, err.Error())
	}
	if strings.Compare(rangeResult, rangeValue) != 0 {
		return logs.LogStep(GetRangeStep, false, "Error, expected value %s not equal to %s from get range.", rangeResult, rangeValue)
	}
	logs.LogStep(GetRangeStep, true, "Success")

	_, err = test.KVHSet(logs, callbacks, hashKey, hashField, hashValue)
	if err != nil {
		return logs.LogStep(SetHashStep, false, "Error setting hash for key: %s. Error was: %s.", hashKey, err.Error())
	}
	defer test.KVDel(logs, callbacks, hashKey)
	logs.LogStep(SetHashStep, true, "Successful.")

	hashResult, err := test.KVHGet(logs, callbacks, hashKey, hashField)
	if err != nil {
		return logs.LogStep(GetHashStep, false, "Error getting hash for key: %s. Error was: %s.", hashKey, err.Error())
	}

	if strings.Compare(hashResult, hashValue) != 0 {
		return logs.LogStep(GetHashStep, false, "Hash returned: %s which was different from expected: %s", hashResult, hashValue)
	}
	logs.LogStep(GetHashStep, true, "Successful. Hash: %s", hashResult)

	hashResults, err := test.KVHGetAll(logs, callbacks, hashKey)
	if err != nil {
		return logs.LogStep(GetAllHashStep, false, "Error, failed to get all hash entries with error: %s", err.Error())
	}

	data := fmt.Sprintf("%v", hashResults[hashField].Row[0])
	if strings.Compare(data, hashValue) != 0 {
		return logs.LogStep(GetAllHashStep, false, "Error, hash field result (%s) and expected value (%s) are not equal.", data, hashValue)
	}
	logs.LogStep(GetAllHashStep, true, "Success")

	return nil
}

func (test *NotebookKVCoreTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	logs.LogInfo("Clean up %s", kvSetKey)
	if _, err := test.KVGet(logs, callbacks, kvSetKey); err == nil {
		_, err := test.KVDel(logs, callbacks, kvSetKey)
		if err != nil {
			logs.LogError("Retrieved kvSetKey but failed to clean up. %s", kvSetKey)
		}
	}

	logs.LogInfo("Clean up %s", kvNXSetKey)
	if _, err := test.KVGet(logs, callbacks, kvNXSetKey); err == nil {
		_, err := test.KVDel(logs, callbacks, kvNXSetKey)
		if err != nil {
			logs.LogError("Retrieved kvNXSetKey but failed to clean up. %s", kvNXSetKey)
		}
	}

	logs.LogInfo("Clean up %s", hashKey)
	if _, err := test.KVHGet(logs, callbacks, hashKey, hashField); err == nil {
		_, err := test.KVDel(logs, callbacks, hashKey)
		if err != nil {
			logs.LogError("Retrieved hashKey but failed to clean up. %s", hashKey)
		}
	}

	return nil
}
