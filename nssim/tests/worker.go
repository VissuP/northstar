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
	"runtime/debug"
	"time"

	. "github.com/verizonlabs/northstar/nssim/utils"
)

const (
	// Define the constants used as Test state.
	Ready    = "Ready"
	Running  = "Running"
	Finished = "Finished"
)

const (
	// Define the constants used as Test status.
	Success             = "Success"
	ConfigurationError  = "ConfigurationError"
	InitializationError = "InitializationError"
	ExecutionError      = "ExecutionError"
	CleanupError        = "CleanupError"
	PanicError          = "PanicError"
	TimeoutError        = "Timeout"
)

// Defines the Scenario Worker status.
type Status struct {
	State          string
	Status         string
	Latency        time.Duration
	NumberOfErrors int
	LastError      string
	Logs           *Logger
}

// Defines the type that implements the service_master.Worker interface
// for test workers.
type TestWorker struct {
	testType    string
	name        string
	repetitions int
	status      chan Status
	logs        *Logger
	callbacks   *ThreadSafeMap
}

// Returns a new test worker
func NewTestWorker(testType string, name string, repetitions int, status chan Status, logs *Logger, callbacks *ThreadSafeMap) *TestWorker {
	logs.LogInfo("NewTestWorker: name:%s", name)

	return &TestWorker{
		testType:    testType,
		name:        name,
		repetitions: repetitions,
		status:      status,
		logs:        logs,
		callbacks:   callbacks,
	}
}

// Runs the test.
func (testWorker *TestWorker) Run(workRoutine int) (err error) {
	testWorker.logs.LogDebug("Starting test worker. id: %d, testType: %s, name: %s, repetitions: %d ",
		workRoutine, testWorker.testType, testWorker.name, testWorker.repetitions)

	startTime := time.Now()
	status := Status{
		State:  Finished,
		Status: Success,
	}

	// Make sure we handle potential panics returned during test initialization,
	// execution, or cleanup. This will make sure we report proper error.
	defer func(status *Status, testWorker *TestWorker) {
		if r := recover(); r != nil {
			testWorker.logs.LogError("Test failed with panic: %v", r)
			testWorker.logs.LogError("Stacktrace: %s", debug.Stack())
			status.State = Finished
			status.Status = PanicError
			status.NumberOfErrors++
			status.LastError = fmt.Errorf("Test %s failed with panic: %v", testWorker.name, r).Error()
			status.Logs = testWorker.logs
			testWorker.status <- *status
		}
	}(&status, testWorker)

	// Get test by resource type
	resourceScenario, err := GetTest(TestId(testWorker.name))

	if err != nil {
		testWorker.logs.LogError("Test with name %s was not found.", testWorker.name)
		status.State = Finished
		status.Status = ConfigurationError
		status.Latency = time.Since(startTime)
		status.LastError = err.Error()
		status.Logs = testWorker.logs
		status.NumberOfErrors++
		status.LastError = err.Error()
		testWorker.status <- status

		return err
	}

	testWorker.logs.LogInfo("Initializing test %s", testWorker.name)

	// Initialize the test. In case of error, return rigth away.
	if err = resourceScenario.Initialize(testWorker.logs); err != nil {
		testWorker.logs.LogError("Test initialization failed with error: %v", err)
		status.State = Finished
		status.Status = InitializationError
		status.Latency = time.Since(startTime)
		status.LastError = err.Error()
		status.Logs = testWorker.logs
		status.NumberOfErrors++
		status.LastError = err.Error()
		testWorker.status <- status

		return err
	}

	testWorker.logs.LogInfo("Executing test %s", testWorker.name)

	for repetition := 0; repetition < testWorker.repetitions; repetition++ {
		if err := resourceScenario.Execute(testWorker.logs, testWorker.callbacks); err != nil {
			testWorker.logs.LogError("Error, worker %d failed to execute test %s with error: %s", workRoutine, testWorker.name, err.Error())
			status.Status = ExecutionError
			status.NumberOfErrors++
			status.LastError = err.Error()
			status.Logs = testWorker.logs
			status.NumberOfErrors++
			status.LastError = err.Error()
		}
	}

	testWorker.logs.LogInfo("Cleaning test %s", testWorker.name)

	if err = resourceScenario.Cleanup(testWorker.logs); err != nil {
		testWorker.logs.LogError("Test cleanup failed with error: %v", err)
		status.Status = CleanupError
		status.LastError = err.Error()
		status.Logs = testWorker.logs
		status.NumberOfErrors++
		status.LastError = err.Error()
	}

	testWorker.logs.LogInfo("Cleaning residual test resources.")
	if err = resourceScenario.CleanupHangingResource(testWorker.logs, testWorker.callbacks); err != nil {
		testWorker.logs.LogError("Failed to cleanup hanging resources with error: %v", err)
	}

	status.Latency = time.Since(startTime)
	testWorker.logs.LogInfo("Completed Test for Worker Index %d: %s-%s, Latency: %s", workRoutine, testWorker.testType, testWorker.name, status.Latency.String())
	testWorker.status <- status

	return nil
}
