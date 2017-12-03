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

package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/pborman/uuid"

	"encoding/base64"
	mrand "math/rand"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
	"github.com/verizonlabs/northstar/pkg/thingspace/models/ott"
)

func VerifyDeviceState(logs *Logger, actual string, expected string) error {
	if actual != expected {
		errString := fmt.Sprintf("Actual device state (%s) does not match expected state (%s)", actual, expected)
		logs.LogError(errString)
		return errors.New(errString)
	}
	return nil
}

func VerifyFieldValue(logs *Logger, actual string, expected string, fieldName string) error {
	if actual != expected {
		errString := fmt.Sprintf("Actual \"%s\" value (%s) does not match expected value (%s)", fieldName, actual, expected)
		logs.LogError(errString)
		return errors.New(errString)
	}
	return nil
}

func VerifyEqual(logs *Logger, actual interface{}, expected interface{}, what string) error {
	if actual != expected {
		errString := fmt.Sprintf("Actual %s value (%v) does not match expected value (%v)", what, actual, expected)
		logs.LogError(errString)
		return errors.New(errString)
	}
	return nil
}

// Helper method used to generate a random numeric id (e.g., IMEI, IMSI)
func GenerateRandomId(length int) string {
	numbers := "012345678901234567890123456789"
	id := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		id[i] = numbers[r.Int()%len(numbers)]
	}

	return string(id)
}

// Helper method used to generate a random String
func GenerateRandomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz"
	id := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		id[i] = chars[r.Int()%len(chars)]
	}

	return string(id)
}

// Returns a unique id that can be used to identify resources.
func GenerateResourceId() string {
	resourceId := uuid.NewRandom()
	return resourceId.String()
}

// find return by matching eventId
func FindEvent(eventId string, events []api.Event) *api.Event {
	for _, event := range events {
		if eventId == event.Id {
			return &event
		}
	}
	return nil
}

func CreateRandomSizeUpdate(logs *Logger) (string, error) {
	size := mrand.Intn(ott.OttUpdateFieldMaxBinaryLength) + 1
	return CreateRandomUpdate(logs, size)
}

func CreateRandomUpdate(logs *Logger, size int) (string, error) {
	logs.LogInfo("creating random update of size %d", size)
	update, err := CreateRandomBytes(logs, size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(update), nil
}

func CreateRandomBytes(logs *Logger, size int) ([]byte, error) {
	logs.LogInfo("creating random byte array of size %d", size)
	bytearray := make([]byte, size)
	if _, err := rand.Read(bytearray); err != nil {
		return nil, fmt.Errorf("Error generating random byte array: %s", err.Error())
	}
	return bytearray, nil
}
