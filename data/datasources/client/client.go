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
	"github.com/verizonlabs/northstar/data/datasources/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/datasources"

type Client interface {
	AddDatasource(accountId string, data *model.DatasourceData) (string, *management.Error)
	DeleteDatasource(accountId string, datasetId string) *management.Error
	GetDatasource(accountId string, datasetId string) (*model.DatasourceData, *management.Error)
	GetDatasources(accountId string) ([]*model.DatasourceData, *management.Error)
	UpdateDatasource(accountId string, datasourceId string, update *model.DatasourceData) *management.Error
}

type DatasourcesClient struct {
	lbClient *lb.LbClient
}

func NewDatasourcesClient() (*DatasourcesClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create datasources data client with error: %s", err.Error())
		return nil, err
	}

	return &DatasourcesClient{lbClient: lbClient}, nil
}

func (client *DatasourcesClient) AddDatasource(accountId string,
	data *model.DatasourceData) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("Datasources data client: Error adding data source: %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *DatasourcesClient) DeleteDatasource(accountId string,
	datasetId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, datasetId)
	err := client.lbClient.Delete(path)
	if err != nil {
		mlog.Error("Datasources data client: Error deleting datasource: %s", err.Error())
		return err
	}
	return nil
}

func (client *DatasourcesClient) GetDatasource(accountId string,
	datasetId string) (*model.DatasourceData, *management.Error) {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, datasetId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Datasources data client: Unable to retrieve datasource: %s", mErr.Error())
		return nil, mErr
	}

	var out *model.DatasourceData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *DatasourcesClient) GetDatasources(accountId string) ([]*model.DatasourceData,
	*management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Datasources data client: Unable to retrieve datasources: %s", mErr.Error())
		return nil, mErr
	}

	var out []*model.DatasourceData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *DatasourcesClient) UpdateDatasource(accountId string,
	datasourceId string,
	update *model.DatasourceData) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, datasourceId)
	_, err := client.lbClient.PutJSON(path, update)
	if err != nil {
		mlog.Error("Datasources data client: Error updating datasource: %s", err.Error())
		return err
	}

	return nil
}
