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
	"github.com/gocql/gocql"
	"time"
)

// Returns the ISO 8601 string representation of the current time.
func GetCurrentTimeInISO8601() string {
	return time.Now().Format(time.RFC3339)
}

// Returns the ISO 8601 string representation of time.
func GetTimeInISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}

// Returns the time representation of the ISO 8601 string.
func GetTimeFromISO8601(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}

// Returns the ISO 8601 string representation of a timeuuid.
func GetTimeUUIDInISO8601(timeuuid string) string {
	// For now, ignore the error.
	uuid, _ := gocql.ParseUUID(timeuuid)

	return GetTimeInISO8601(uuid.Time())
}

// Returns true if the specified datetime string if valid ISO 8601.
func IsValidISO8601(t string) bool {
	_, err := time.Parse(time.RFC3339, t)

	if err != nil {
		return false
	}

	return true
}
