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
	"github.com/verizonlabs/northstar/dpe-stream/master/model"
	"github.com/verizonlabs/northstar/dpe-stream/master/util"
)

const BASE_URI = util.StreamBasePath + "/jobs"

type StreamClient struct {
	lbClient *lb.LbClient
}

func NewStreamClient() (*StreamClient, error) {
	url, err := util.GetStreamBaseUrl()
	if err != nil {
		mlog.Error("Failed to get dpe stream base url with error: %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create client with error: %s", err.Error())
		return nil, err
	}
	return &StreamClient{lbClient: lbClient}, nil
}

func (client *StreamClient) StartJob(accountId string,
	job *model.StreamJob) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, job)
	if err != nil {
		mlog.Error("DPE stream client: Error provisioning: %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *StreamClient) DeleteJob(accountId string, jobId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, jobId)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}
	return nil
}
