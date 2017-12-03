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
	"github.com/verizonlabs/northstar/data/mappings/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/mappings"

type Client interface {
	AddMapping(accountId string, data *model.MappingsData) (string, *management.Error)
	ListMappings(accountId string) ([]*model.MappingsData, *management.Error)
	DeleteMapping(accountId string, mappingId string) *management.Error
	GetMapping(accountId string, mappingId string) (*model.MappingsData, *management.Error)
	GetMappingByEventId(accountId string, eventId string) (*model.MappingsData, *management.Error)
}

type MappingsClient struct {
	lbClient *lb.LbClient
}

func NewMappingsClient() (*MappingsClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create kafka mngr data client with error: %s", err.Error())
		return nil, err
	}

	return &MappingsClient{lbClient: lbClient}, nil
}

func (client *MappingsClient) AddMapping(accountId string,
	data *model.MappingsData) (string, *management.Error) {
	path := fmt.Sprintf("%s/by-accountid/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("Mapping dataservice client: Error adding data %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *MappingsClient) ListMappings(accountId string) ([]*model.MappingsData, *management.Error) {
	path := fmt.Sprintf("%s/by-accountid/%s", BASE_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Mapping dataservice client: Error listing data: %s", mErr.Error())
		return nil, mErr
	}

	var out []*model.MappingsData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *MappingsClient) DeleteMapping(accountId string, mappingId string) *management.Error {
	path := fmt.Sprintf("%s/by-accountid/%s/%s", BASE_URI, accountId, mappingId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}

	return nil
}

func (client *MappingsClient) GetMapping(accountId string,
	mappingId string) (*model.MappingsData, *management.Error) {
	path := fmt.Sprintf("%s/by-accountid/%s/%s", BASE_URI, accountId, mappingId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Mappings dataservice client: Error getting mapping: %s", mErr.Error())
		return nil, mErr
	}

	var mapping *model.MappingsData
	if err := json.Unmarshal(resp, &mapping); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return mapping, nil
}

func (client *MappingsClient) GetMappingByEventId(accountId string,
	eventId string) (*model.MappingsData, *management.Error) {
	path := fmt.Sprintf("%s/by-eventid/%s/%s", BASE_URI, accountId, eventId)
	resp, err := client.lbClient.Get(path)
	if err != nil {
		mlog.Error("Mappings dataservice client: Error getting mapping by event id: %v", resp)
		return nil, err
	}

	var mapping *model.MappingsData
	if err := json.Unmarshal(resp, &mapping); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return mapping, nil
}
