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
	"github.com/verizonlabs/northstar/data/snippets/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/snippets"

type Client interface {
	AddSnippet(accountId string, add *model.SnippetData) (string, *management.Error)
	DeleteSnippet(accountId string, snippetId string) *management.Error
	GetSnippets(accountId string) ([]*model.SnippetData, *management.Error)
	GetSnippet(accountId string, snippetId string) (*model.SnippetData, *management.Error)
	UpdateSnippet(accountId string, snippetId string, update *model.SnippetData) *management.Error
}

type SnippetsClient struct {
	lbClient *lb.LbClient
}

func NewSnippetClient() (*SnippetsClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create snippets data client with error: %s", err.Error())
		return nil, err
	}

	return &SnippetsClient{lbClient: lbClient}, nil
}

func (client *SnippetsClient) AddSnippet(accountId string,
	add *model.SnippetData) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, add)
	if err != nil {
		mlog.Error("Snippets dataservice client: Error adding data: %v", err.Error())
		return "", err
	}

	return string(resp), nil
}

func (client *SnippetsClient) DeleteSnippet(accountId string,
	snippetId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, snippetId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}

	return nil
}

func (client *SnippetsClient) GetSnippets(accountId string) ([]*model.SnippetData,
	*management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Snippets dataservice client: Error listing data: %v", mErr.Error())
		return nil, mErr
	}

	var out []*model.SnippetData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *SnippetsClient) GetSnippet(accountId string,
	snippetId string) (*model.SnippetData, *management.Error) {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, snippetId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Snippets dataservice client: Error getting data %s", mErr.Error())
		return nil, mErr
	}

	var out *model.SnippetData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *SnippetsClient) UpdateSnippet(accountId string,
	snippetId string,
	update *model.SnippetData) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, snippetId)
	resp, err := client.lbClient.PutJSON(path, update)
	if err != nil {
		mlog.Error("Snippets dataservice client: Error updating code %v", resp)
		return err
	}

	return nil
}
