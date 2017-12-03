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
	"github.com/verizonlabs/northstar/cron/model"
	"github.com/verizonlabs/northstar/cron/util"
)

const BASE_URI = util.CronBasePath + "/jobs"

type CronClient struct {
	lbClient *lb.LbClient
}

func NewCronClient() (*CronClient, error) {
	url, err := util.GetCronBaseUrl()
	if err != nil {
		mlog.Error("Failed to get cron base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create cron client with error: %s", err.Error())
		return nil, err
	}

	return &CronClient{lbClient: lbClient}, nil
}

func (client *CronClient) AddJob(accountId string, job *model.Job) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, job)
	if err != nil {
		mlog.Error("Cron client: Error adding: %v", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *CronClient) DeleteJob(accountId string, jobId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, jobId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}
	return nil
}

func (client *CronClient) UpdateJob(accountId string,
	snippetId string,
	update *model.Job) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, snippetId)
	_, err := client.lbClient.PutJSON(path, update)
	if err != nil {
		mlog.Error("Cron client: Error updating: %s", err.Error())
		return err
	}
	return nil
}
