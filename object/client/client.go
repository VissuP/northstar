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
	"github.com/verizonlabs/northstar/object/model"
	"github.com/verizonlabs/northstar/object/util"
)

const (
	BUCKETS_URI = util.ObjectBasePath + "/buckets"
	FILES_URI   = util.ObjectBasePath + "/files"
)

type ObjectClient struct {
	lbClient *lb.LbClient
}

func NewObjectClient() (*ObjectClient, error) {
	url, err := util.GetObjectBaseUrl()
	if err != nil {
		mlog.Error("Failed to get object base url with error: %s", err.Error())
		return nil, err
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create object client with error: %s", err.Error())
		return nil, err
	}

	return &ObjectClient{lbClient: lbClient}, nil
}

func (client *ObjectClient) CreateBucket(accountId string,
	bucket *model.Bucket) (string, *management.Error) {
	path := fmt.Sprintf("%s/%s", BUCKETS_URI, accountId)
	resp, err := client.lbClient.PostJSON(path, bucket)
	if err != nil {
		mlog.Error("Object client: Error creating bucket: %v", err.Error())
		return "", err
	}
	return string(resp), nil
}

func (client *ObjectClient) DeleteBucket(accountId string, bucketName string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BUCKETS_URI, accountId, bucketName)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}
	return nil
}

func (client *ObjectClient) ListBuckets(accountId string) ([]model.Bucket, *management.Error) {
	path := fmt.Sprintf("%s/%s", BUCKETS_URI, accountId)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Object client: Error listing buckets: %v", mErr.Error())
		return nil, mErr
	}

	var out []model.Bucket
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *ObjectClient) UploadFile(accountId string,
	bucketName string,
	data *model.UploadData) (string, *management.Error) {

	path := fmt.Sprintf("%s/%s/%s", FILES_URI, accountId, bucketName)
	mlog.Debug("Upload path: %s", path)

	resp, mErr := client.lbClient.PostJSON(path, data)
	if mErr != nil {
		return "", mErr
	}

	return string(resp), nil
}

func (client *ObjectClient) DownloadFile(accountId,
	bucketName,
	fileName string) (*model.DownloadData, *management.Error) {

	path := fmt.Sprintf("%s/%s/%s/%s", FILES_URI, accountId, bucketName, fileName)
	mlog.Debug("Download path: %s", path)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Object client: Error listing buckets: %v", mErr.Error())
		return nil, mErr
	}

	var data *model.DownloadData
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return data, nil
}

func (client *ObjectClient) DeleteFile(accountId, bucketName, fileName string) *management.Error {
	path := fmt.Sprintf("%s/%s/%s/%s", FILES_URI, accountId, bucketName, fileName)
	err := client.lbClient.Delete(path)
	if err != nil {
		return err
	}
	return nil
}

func (client *ObjectClient) ListFiles(accountId,
	bucketName string) ([]model.Object, *management.Error) {
	path := fmt.Sprintf("%s/%s/%s", FILES_URI, accountId, bucketName)
	resp, mErr := client.lbClient.Get(path)
	if mErr != nil {
		mlog.Error("Object client: Error listing buckets: %v", mErr.Error())
		return nil, mErr
	}

	var out []model.Object
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}
