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
	"encoding/json"
	"fmt"
	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/events/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/events"

type Client interface {
	AddEvent(accountId string, data *model.EventData) (string, *management.Error)
	DeleteEvent(accountId string, eventId string) *management.Error
	ListEvents(accountId string) ([]*model.EventData, *management.Error)
}

type EventsClient struct {
	lbClient *lb.LbClient
}

func NewEventsClient() (*EventsClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create events data client with error: %s", err.Error())
		return nil, err
	}

	return &EventsClient{lbClient: lbClient}, nil
}

func (client *EventsClient) AddEvent(accountId string,
	data *model.EventData) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("Events data client: Error adding event: %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *EventsClient) DeleteEvent(accountId string, eventId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, eventId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}

	return nil
}

func (client *EventsClient) ListEvents(accountId string) ([]*model.EventData, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Events data client: Error getting listing events: %v", mErr.Error())
		return nil, mErr
	}

	var out []*model.EventData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}
