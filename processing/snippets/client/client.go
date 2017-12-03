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
	"github.com/verizonlabs/northstar/processing/snippets/model"
	"github.com/verizonlabs/northstar/processing/util"
)

const BASE_URI = util.ProcessingBasePath + "/snippets"

type Client interface {
	StartSnippet(accountId string, snippet *model.Snippet) (string, *management.Error)
	StopSnippet(accountId string, invocationId string) *management.Error
}

type SnippetsClient struct {
	lbClient *lb.LbClient
}

func NewSnippetsClient() (*SnippetsClient, error) {
	url, err := util.GetProcessingBaseUrl()
	if err != nil {
		mlog.Error("Failed to get snippets base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create snippet processing client with error: %s", err.Error())
		return nil, err
	}

	return &SnippetsClient{lbClient: lbClient}, nil
}

func (client *SnippetsClient) StartSnippet(accountId string,
	snippet *model.Snippet) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, snippet)
	if err != nil {
		mlog.Error("Snippets processing client: Error starting: %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *SnippetsClient) StopSnippet(accountId string,
	invocationId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, invocationId)
	err := client.lbClient.Delete(path)
	if err != nil {
		mlog.Error("Snippets processing client: Error stopping %s", err.Error())
		return err
	}
	return nil
}
