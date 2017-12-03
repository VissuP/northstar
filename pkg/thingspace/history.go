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

package thingspace

import (
	"encoding/json"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/thingspace/api"
)


// Returns the field history that meet the specified filter.
func (client *ThingSpaceClient) QueryFieldHistory(filter interface{}) ([]api.Event, *management.Error) {
	mlog.Debug("QueryFieldHistory: filter:%#v", filter)

	url := client.baseUrl + api.BASE_PATH + "events"
	body, err := client.executeHttpMethod("GET", url, "", filter)

	if err != nil {
		mlog.Error("Error, failed to get field history with error: %s", url, err.Description)
		return nil, err
	}

	events := make([]api.Event, 0)
	if err := json.Unmarshal(body, &events); err != nil {
		mlog.Error("Error, failed to unmarshal response to events with error: %s.", err.Error())
		return nil, management.ErrorInternal
	}

	return events, nil
}
