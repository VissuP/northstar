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

package client

import (
	"fmt"
	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/processing/events/model"
	"github.com/verizonlabs/northstar/processing/util"
)

const BASE_URI = util.ProcessingBasePath + "/events"

type EventsClient struct {
	lbClient *lb.LbClient
}

func NewEventsClient() (*EventsClient, error) {
	url, err := util.GetProcessingBaseUrl()
	if err != nil {
		mlog.Error("Failed to get events base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create snippet processing client with error: %s", err.Error())
		return nil, err
	}

	return &EventsClient{lbClient: lbClient}, nil
}

func (client *EventsClient) InvokeEvent(accountId string,
	eventId string,
	options *model.Options) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, eventId)
	resp, err := client.lbClient.PostJSON(path, options)
	if err != nil {
		mlog.Error("Events user client: Error invoke: %s", err.Error())
		return "", err
	}

	return string(resp), nil
}
