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
	"github.com/verizonlabs/northstar/data/datasets/model"
	"github.com/verizonlabs/northstar/data/util"
)

const BASE_URI = util.DataBasePath + "/datasets"

type Client interface {
	AddDataset(accountId string, data *model.DatasetData) (string, *management.Error)
	DeleteDataset(accountId string, datasetId string) *management.Error
	GetDatasetById(accountId string, datasetId string) (*model.DatasetData, *management.Error)
	GetDatasetByName(accountId string, name string) (*model.DatasetData, *management.Error)
	GetDatasets(accountId string) ([]*model.DatasetData, *management.Error)
	UpdateDataset(accountId string, datasetId string, update *model.DatasetData) *management.Error
}

type DatasetsClient struct {
	lbClient *lb.LbClient
}

func NewDatasetsClient() (*DatasetsClient, error) {
	url, err := util.GetDataBaseUrl()
	if err != nil {
		mlog.Error("Failed to get data base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create datasets data client with error: %s", err.Error())
		return nil, err
	}

	return &DatasetsClient{lbClient: lbClient}, nil
}

func (client *DatasetsClient) AddDataset(accountId string,
	data *model.DatasetData) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, data)
	if err != nil {
		mlog.Error("Datasets data client: Error adding dataset: %s", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *DatasetsClient) DeleteDataset(accountId string, datasetId string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, datasetId)
	err := client.lbClient.Delete(path)
	if err != nil {
		mlog.Error("Datasets data client: Error deleting dataset: %s", err.Error())
		return err
	}
	return nil
}

func (client *DatasetsClient) GetDatasetById(accountId string,
	datasetId string) (*model.DatasetData, *management.Error) {
	path := fmt.Sprintf("%s/%s/by-id/%s", BASE_URI, accountId, datasetId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Datasets data client: Unable to retrieve dataset: %s", mErr.Error())
		return nil, mErr
	}

	var out *model.DatasetData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *DatasetsClient) GetDatasetByName(accountId string,
	name string) (*model.DatasetData, *management.Error) {
	path := fmt.Sprintf("%s/%s/by-name/%s", BASE_URI, accountId, name)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Datasets data client: Unable to retrieve dataset: %s", mErr.Error())
		return nil, mErr
	}

	var out *model.DatasetData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *DatasetsClient) GetDatasets(accountId string) ([]*model.DatasetData, *management.Error) {
	path := fmt.Sprintf("%s/%s", BASE_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Datasets data client: Unable to retrieve datasets: %s", mErr.Error())
		return nil, mErr
	}

	var out []*model.DatasetData
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *DatasetsClient) UpdateDataset(accountId string,
	datasetId string,
	update *model.DatasetData) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, accountId, datasetId)
	_, err := client.lbClient.PutJSON(path, update)
	if err != nil {
		mlog.Error("Datasets data client: Error updating dataset: %s", err.Error())
		return err
	}

	return nil
}
