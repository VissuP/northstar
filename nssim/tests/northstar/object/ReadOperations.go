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

package execution

import (
	"encoding/json"
	"fmt"
	"time"

	"bytes"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.ObjectRead), NewObjectReadTest)
}

const (
	executionTimeout = 10 * time.Second
)

//ObjectReadTest is the implementation of the test interface for Object Read functionality.
type ObjectReadTest struct {
	*northstar.NorthstarApiBaseTest
}

//NewObjectReadTest creates an instance of the Object Read Test
func NewObjectReadTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &ObjectReadTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//steps
	createBucketStep utils.Step = "Create bucket"
	listBucketsStep  utils.Step = "List buckets."
	createObjectStep utils.Step = "Create Object"
	listObjectsStep  utils.Step = "List objects"
	getObjectStep    utils.Step = "Get Object"

	//common variables
	bucketName = "test-bucket"
	objectName = "test-file-from-string"
)

//Execute executes the test
func (test *ObjectReadTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- ObjectReadTest")

	//Create the bucket for our test
	err := test.createBucket(logs, bucketName, callbacks)
	if err != nil {
		return logs.LogStep(createBucketStep, false, "Create bucket failed. %s", err.Error())
	}
	defer test.deleteBucket(logs, bucketName, callbacks)
	logs.LogStep(createBucketStep, true, "Create bucket succeeded.")

	buckets, mErr := test.NorthstarApiClient.ListBuckets(test.Users[0].Token.AccessToken)
	if mErr != nil {
		return logs.LogStep(listBucketsStep, false, "Failed to list buckets with error: %s", mErr.Description)
	}

	if len(buckets) == 0 {
		return logs.LogStep(listBucketsStep, false, "Error, bucket list was empty.")
	}
	logs.LogStep(listBucketsStep, true, "Success. Buckets: %+v", buckets)

	expectedData := &northstarApiModel.Data{
		ContentType: "text/plain",
		Payload:     []byte("test-data"),
	}

	err = test.createObject(logs, bucketName, objectName, expectedData, callbacks)
	if err != nil {
		return logs.LogStep(createObjectStep, false, "Create object failed. %s", err.Error())
	}
	defer test.deleteObject(logs, bucketName, objectName, callbacks)

	objects, mErr := test.NorthstarApiClient.ListObjects(test.Users[0].Token.AccessToken, bucketName, "", 1000, "")
	if mErr != nil {
		return logs.LogStep(listObjectsStep, false, "Failed to list objects with error: %s", mErr.Description)
	}

	if !test.objectExists(logs, objects, objectName) {
		return logs.LogStep(listObjectsStep, false, "Failed to find our new object in listed objects.")
	}
	logs.LogStep(listObjectsStep, true, "Success. Objects: %+v", objects)

	object, mErr := test.NorthstarApiClient.GetObject(test.Users[0].Token.AccessToken, bucketName, objectName)
	if mErr != nil {
		return logs.LogStep(getObjectStep, false, "Failed to get object with error: %s", mErr.Description)
	}
	err = test.verifyObject(logs, expectedData, object)
	if err != nil {
		return logs.LogStep(getObjectStep, false, "Failed to validate object. %s", err.Error())
	}
	logs.LogStep(getObjectStep, true, "Successfully get object. Object: %+v", object)

	return nil
}

//createObject creates the specified object
func (test *ObjectReadTest) createObject(logs *utils.Logger, bucket string, object string, data *northstarApiModel.Data, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("Create an object")

	logs.LogInfo("Sending request to create object.")
	code := `local object = require("nsObject")

					function main()
    					local err = object.uploadFile("` + bucket + `","` + object + `", "` + string(data.Payload) + `", "` + data.ContentType + `")
    					if err ~= nil then
        					error(err)
    					end
					end`


	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return logs.LogError("Failed to create object. %s", err.Error())
	}

	output := northstarApiModel.Output{}
	err = json.Unmarshal([]byte(*response), &output)
	if err != nil {
		return err
	}

	logs.LogInfo("Execution output: %+v results: %+v", output, output.ExecutionResults)

	if output.Status != northstarApiModel.OutputSuccessStatus {
		return logs.LogError("Failed to execute with error: %s.", output.StatusDescription)
	}

	return nil
}

//deleteObject deletes the specified object
func (test *ObjectReadTest) deleteObject(logs *utils.Logger, bucket string, object string, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("Delete object")

	code := `local object = require("nsObject")

					function main()
    					local err = object.deleteFile("` + bucket + `", "` + object + `")
    					if err ~= nil then
        					error(err)
    					end
					end`

	logs.LogInfo("Sending request to delete object")
	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return logs.LogError("Failed to create object. %s", err.Error())
	}

	output := northstarApiModel.Output{}
	err = json.Unmarshal([]byte(*response), &output)
	if err != nil {
		return err
	}

	logs.LogInfo("Execution output: %+v results: %+v", output, output.ExecutionResults)

	if output.Status != northstarApiModel.OutputSuccessStatus {
		return logs.LogError("Failed to execute with error: %s.", output.StatusDescription)
	}


	return nil
}

//CreateBucket creates the specified bucket
func (test *ObjectReadTest) createBucket(logs *utils.Logger, bucket string, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("Create a bucket")


	code := `local object = require("nsObject")

					function main()
    					local err = object.createBucket("` + bucket + `")
    					if err ~= nil then
        					error(err)
    					end
					end`

	logs.LogInfo("Sending request to create bucket.")
	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return logs.LogError("Failed to create object. %s", err.Error())
	}

	output := northstarApiModel.Output{}
	err = json.Unmarshal([]byte(*response), &output)
	if err != nil {
		return err
	}

	logs.LogInfo("Execution output: %+v results: %+v", output, output.ExecutionResults)

	if output.Status != northstarApiModel.OutputSuccessStatus {
		return logs.LogError("Failed to execute with error: %s.", output.StatusDescription)
	}

	return nil
}

//DeleteBucket deletes the specified bucket
func (test *ObjectReadTest) deleteBucket(logs *utils.Logger, bucket string, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("Cleanup Bucket")


	code := `local object = require("nsObject")

				function main()
    				local err = object.deleteBucket("` + bucket + `")
    				if err ~= nil then
        				error(err)
    				end
				end`

	logs.LogInfo("Sending request to delete bucket.")
	response, err := test.ExecuteCell(logs, callbacks, code)
	if err != nil {
		return logs.LogError("Failed to create object. %s", err.Error())
	}

	output := northstarApiModel.Output{}
	err = json.Unmarshal([]byte(*response), &output)
	if err != nil {
		return err
	}

	logs.LogInfo("Execution output: %+v results: %+v", output, output.ExecutionResults)

	if output.Status != northstarApiModel.OutputSuccessStatus {
		return logs.LogError("Failed to execute with error: %s.", output.StatusDescription)
	}

	return nil
}

//verifyObject verifies that two objects are equal
func (test *ObjectReadTest) verifyObject(logs *utils.Logger, expectedData *northstarApiModel.Data, actualData *northstarApiModel.Data) error {
	logs.LogInfo("verifyObject")
	if expectedData.ContentType != actualData.ContentType {
		return fmt.Errorf("Error, expected object content type (%s) does not match created content type (%s).", expectedData.ContentType, actualData.ContentType)
	}

	if bytes.Compare(expectedData.Payload, actualData.Payload) != 0 {
		return fmt.Errorf("Error, expected object content (%s) does not match created content (%s).", expectedData.Payload, actualData.Payload)
	}
	return nil
}

//objectExists confirms that an object exists in the collection
func (test *ObjectReadTest) objectExists(logs *utils.Logger, objects []northstarApiModel.Object, objectName string) bool {
	logs.LogInfo("Objects: %+v", objects)

	for _, object := range objects {
		if object.Key == objectName {
			logs.LogInfo("Found object in list")
			return true
		}
	}
	return false
}

func (test *ObjectReadTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	//Check to see if we can get our bucket.
	objects, err := test.NorthstarApiClient.ListObjects(test.Users[0].Token.AccessToken, bucketName, "", 1000, "")
	if err != nil {
		return logs.LogError("Failed to list objects. %s", err.Error())
	}

	for _, object := range objects {
		if err := test.deleteObject(logs, bucketName, object.Key, callbacks); err != nil {
			logs.LogError("Failed to delete object. %s", err.Error())
			continue
		}
	}

	return nil
}
