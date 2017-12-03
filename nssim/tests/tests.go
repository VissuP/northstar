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

package tests

import (
	"fmt"

	"github.com/verizonlabs/northstar/nssim/utils"
)

var (
	tests = make(map[TestId]TestCreator)
)

// Define the test interface
type Test interface {
	Initialize(logs *utils.Logger) error
	Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error
	Cleanup(logs *utils.Logger) error
	CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error
}

// Defines a test id type
type TestId string

// Defines a test creator (i.e., factory method).
type TestCreator func() (Test, error)

// Register a store creator for the specified id.
func Register(id TestId, creator TestCreator) {
	tests[id] = creator
}

// Returns the store for the specified id.
func GetTest(id TestId) (Test, error) {
	creator, found := tests[id]
	if found == false {
		return nil, fmt.Errorf("Error, test with id %s not found.", id)
	}

	test, err := creator()
	if err != nil {
		return nil, fmt.Errorf("Error, could not create test with error: %s", err.Error())
	}

	return test, nil
}
