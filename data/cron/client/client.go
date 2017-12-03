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
	"github.com/verizonlabs/northstar/data/cron/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/cron"

type Client interface {
	AddJob(accountId string, data *model.JobData) *management.Error
	DeleteJob(accountId string, jobId string) *management.Error
	GetJob(accountId string, jobId string) (*model.JobData, *management.Error)
	GetAllJobs() ([]*model.JobData, *management.Error)
	GetJobsByAccountId(accountId string) ([]*model.JobData, *management.Error)
	UpdateJob(accountId string, jobId string, update *model.JobData) *management.Error
}

type CronClient struct {
	lbClient *lb.LbClient
}

func NewCronClient() (*CronClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create cron data client for url with error: %s", err.Error())
		return nil, err
	}

	return &CronClient{lbClient: lbClient}, nil
}

func (client *CronClient) AddJob(accountId string, data *model.JobData) *management.Error {
	path := fmt.Sprintf("%s/by-accountid/%s", BASE_URI, accountId)
	_, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("Cron data client: Error adding job: %v", err.Error())
		return err
	}
	return nil
}

func (client *CronClient) DeleteJob(accountId string, jobId string) *management.Error {
	path := fmt.Sprintf("%s/by-accountid/%s/%s", BASE_URI, accountId, jobId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}

	return nil
}

func (client *CronClient) GetJob(accountId string,
	jobId string) (*model.JobData, *management.Error) {
	path := fmt.Sprintf("%s/by-accountid/%s/%s", BASE_URI, accountId, jobId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Cron data client: Error getting job: %v", mErr.Error())
		return nil, mErr
	}

	var out *model.JobData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *CronClient) GetAllJobs() ([]*model.JobData, *management.Error) {
	path := fmt.Sprintf("%s/all/jobs", BASE_URI)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Cron data client: Error listing all jobs: %v", mErr.Error())
		return nil, mErr
	}

	var out []*model.JobData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *CronClient) GetJobsByAccountId(accountId string) ([]*model.JobData,
	*management.Error) {
	path := fmt.Sprintf("%s/by-accountid/%s", BASE_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Cron data client: Error listing jobs by account id: %v", mErr.Error())
		return nil, mErr
	}

	var out []*model.JobData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *CronClient) UpdateJob(accountId string,
	jobId string,
	update *model.JobData) *management.Error {
	path := fmt.Sprintf("%s/by-accountid/%s/%s", BASE_URI, accountId, jobId)
	_, err := client.lbClient.PutJSON(path, update)
	if err != nil {
		mlog.Error("Cron data client: Error updating job: %s", err.Error())
		return err
	}

	return nil
}
