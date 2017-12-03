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

package internalClient

import (
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

const (
	ExecutionsPath = "callbacks"
)

// ExecutionCallback processess the callback for an execution.
func (client *InternalClient) ExecutionCallback(response *model.ExecutionResponse) *management.Error {
	mlog.Debug("ExecutionCallback: response:%+v", response)
	path := client.getResourcePath(ExecutionsPath) + "/execution"

	// If error, return.
	if _, mErr := client.lbClient.PostJSON(path, response); mErr != nil {
		return mErr
	}

	return nil
}
