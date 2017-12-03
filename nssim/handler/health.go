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

package handler

import (
	"net/http"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/nssim/config"
)

var (
	serviceStates = NewStates()
)

const (
	S_HEALTH string = "healthy"
)

// States defines the default service states map
type States struct {
	m map[string]string
}

// GetHealth returns the Portal Service health.
func (controller *Controller) GetHealth() *management.Health {
	mlog.Debug("GetHealth")

	health := &management.Health{
		HttpStatus:  http.StatusOK,
		Name:        config.Configuration.ServiceName,
		Id:          S_HEALTH,
		Description: serviceStates.Get(S_HEALTH),
	}

	return health
}

// NewStates returns default state.
func NewStates() *States {
	states := &States{m: make(map[string]string)}
	states.m[S_HEALTH] = "Service running stable."
	return states
}

// Get returns the description of the state with the specified id.
func (states *States) Get(id string) string {
	if m, ok := states.m[id]; ok {
		return m
	}

	return id
}
