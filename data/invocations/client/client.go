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
	"github.com/verizonlabs/northstar/data/invocations/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/invocations"

type Client interface {
	AddInvocation(accountId string, data *model.InvocationData) (string, *management.Error)
	UpdateInvocation(accountId string, invocationId string, output *model.InvocationData) *management.Error
	GetInvocation(accountId string, invocationId string) (*model.InvocationData, *management.Error)
	GetInvocationsByAccountId(accountId string, limit int) ([]*model.InvocationData, *management.Error)
	GetInvocationResults(accountId string, snippetId string, limit int) ([]*model.InvocationData, *management.Error)
	DeleteInvocation(accountId string, invocationId string) *management.Error
}

type InvocationClient struct {
	lbClient *lb.LbClient
}

func NewInvocationClient() (*InvocationClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create invocation data client with error: %s", err.Error())
		return nil, err
	}

	return &InvocationClient{lbClient: lbClient}, nil
}

func (client *InvocationClient) AddInvocation(accountId string,
	data *model.InvocationData) (string, *management.Error) {
	path := fmt.Sprintf("%s/invocation/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("Invocation dataservice client: Error adding data: %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *InvocationClient) UpdateInvocation(accountId string,
	invocationId string,
	output *model.InvocationData) *management.Error {
	path := fmt.Sprintf("%s/invocation/%s/%s", BASE_URI, accountId, invocationId)
	_, err := client.lbClient.PostJSON(path, output)
	if err != nil {
		mlog.Error("Invocation dataservice client: Error storing data: %s", err.Error())
		return err
	}
	return nil
}

func (client *InvocationClient) GetInvocation(accountId string,
	invocationId string) (*model.InvocationData, *management.Error) {
	path := fmt.Sprintf("%s/invocation/%s/%s", BASE_URI, accountId, invocationId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Invocation dataservice client: Error getting invocation: %s", mErr.Error())
		return nil, mErr
	}

	var out *model.InvocationData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *InvocationClient) GetInvocationsByAccountId(accountId string,
	limit int) ([]*model.InvocationData,
	*management.Error) {
	path := fmt.Sprintf("%s/history/by-account/%s/%d", BASE_URI, accountId, limit)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Invocation dataservice client: Error getting invocations by account id: %s",
			mErr.Error())
		return nil, mErr
	}

	var out []*model.InvocationData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *InvocationClient) GetInvocationResults(accountId string,
	snippetId string,
	limit int) ([]*model.InvocationData, *management.Error) {
	path := fmt.Sprintf("%s/history/by-snippet/%s/%s/%d", BASE_URI, accountId, snippetId, limit)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		return nil, mErr
	}

	var results []*model.InvocationData
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return results, nil
}

func (client *InvocationClient) DeleteInvocation(accountId string,
	invocationId string) *management.Error {
	path := fmt.Sprintf("%s/invocation/%s/%s", BASE_URI, accountId, invocationId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}

	return nil
}
