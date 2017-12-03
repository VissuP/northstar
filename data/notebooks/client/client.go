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
	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/notebooks/model"
	"github.com/verizonlabs/northstar/data/util"
)

type Client interface {
	CreateAccess(access *model.Access) (*model.Access, *management.Error)
	UpdateAccess(accessId string, access *model.Access) (*model.Access, *management.Error)
	QueryAccess(query *model.Query) ([]model.Access, *management.Error)
	DeleteAccess(accessId string) *management.Error
	CreateNotebook(notebook *model.Notebook) (*model.Notebook, *management.Error)
	UpdateNotebook(notebookId string, notebook *model.Notebook) (*model.Notebook, *management.Error)
	GetNotebook(notebookId string) (*model.Notebook, *management.Error)
	DeleteNotebook(notebookId string) *management.Error
}

type NotebooksClient struct {
	baseUrl  string
	lbClient *lb.LbClient
}

// Returns a new notebooks client.
func NewNotebooksClient() (*NotebooksClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create notebook data client with error: %s", err.Error())
		return nil, err
	}

	return &NotebooksClient{lbClient: lbClient}, nil
}
