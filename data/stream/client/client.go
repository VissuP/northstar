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
	"net/http"
	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/stream/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/stream"

type StreamClient struct {
	lbClient *lb.LbClient
}

type Client interface {
	AddJob(accountId string, data *model.JobData) *management.Error
	GetJob(accountId string, jobId string) (*model.JobData, *management.Error)
	GetJobs(accountId string) ([]*model.JobData, *management.Error)
	UpdateJob(accountId string, jobId string, update *model.JobData) *management.Error
	DeleteJob(accountId string, jobId string) *management.Error
}

func NewStreamClient() (*StreamClient, error) {
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

	return &StreamClient{lbClient: lbClient}, nil
}

func (client *StreamClient) AddJob(accountId string, data *model.JobData) *management.Error {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	_, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("DPE stream data client: Error adding job: %s", err.Error())
		return err
	}
	return nil
}

func (client *StreamClient) GetJob(accountId string, jobId string) (*model.JobData, *management.Error) {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, jobId)
	resp, err := client.lbClient.Get(path)
	if err != nil {
		mlog.Alarm("DPE stream data client: Error retrieving job: %s", err.Error())
		return nil, err
	}

	var out *model.JobData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.NewError(http.StatusInternalServerError, "server_error", err.Error())
	}

	return out, nil
}

func (client *StreamClient) GetJobs(accountId string) ([]*model.JobData, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.Get(path)
	if err != nil {
		mlog.Error("DPE stream data client: Error listing jobs: %s", err.Error())
		return nil, err
	}

	var out []*model.JobData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.NewError(http.StatusInternalServerError, "server_error", err.Error())
	}

	return out, nil
}

func (client *StreamClient) UpdateJob(accountId string,
	jobId string,
	update *model.JobData) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, jobId)
	_, err := client.lbClient.PutJSON(path, update)
	if err != nil {
		mlog.Error("DPE stream data client: Error updating job: %s", err.Error())
		return err
	}

	return nil
}

func (client *StreamClient) DeleteJob(accountId string, jobId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, jobId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}

	return nil
}
