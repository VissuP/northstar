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

package scheduler

import (
	"fmt"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

var (
	schedulers = make(map[string]Scheduler)
)

// Defines the interface used to support transformation schedules.
type Scheduler interface {
	// Define common resource operations.
	Create(accountId string, transformationId string, schedule *model.Schedule) (*model.Schedule, *management.Error)
	Get(accountId string, scheduleId string) (*model.Schedule, *management.Error)
	Delete(accountId string, scheduleId string) *management.Error
}

// Register a scheduler for a specific id.
func Register(id string, provider Scheduler) {
	schedulers[id] = provider
}

// Returns the schduler for the specified id.
func Get(id string) (Scheduler, error) {
	provider, found := schedulers[id]

	if found == false {
		return nil, fmt.Errorf("Scheduler for id %s not found.", id)
	}

	return provider, nil
}
