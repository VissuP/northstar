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
	"strconv"

	"net/http"

	"io/ioutil"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

const (
	//ObjectsPath is the endpoint on the northstarservice for objects
	ObjectsPath = "objects"
)

//ListBuckets lists the buckets
func (client *Client) ListBuckets(accessToken string) ([]model.Bucket, *management.Error) {
	mlog.Debug("ListBuckets")

	path := client.getResourcePath(ObjectsPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)
	if mErr != nil {
		return nil, mErr
	}

	buckets := []model.Bucket{}
	if err := json.Unmarshal(response, &buckets); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to unmarshal body with error: %s", err.Error()))
	}

	return buckets, nil
}

//ListObjects lists the objects with the subpath
func (client *Client) ListObjects(accessToken string, bucket string, prefix string, count int, marker string) ([]model.Object, *management.Error) {
	mlog.Debug("ListObjects")

	path := client.getResourcePath(ObjectsPath) + "/" + bucket + "/list/" + prefix + "?count=" + strconv.Itoa(count) + "&marker=" + marker
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)
	if mErr != nil {
		return nil, mErr
	}

	objects := []model.Object{}
	if err := json.Unmarshal(response, &objects); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to unmarshal body with error: %s", err.Error()))
	}

	return objects, nil

}

//GetObject returns the specified object
func (client *Client) GetObject(accessToken string, bucket string, object string) (*model.Data, *management.Error) {
	mlog.Debug("GetObject")

	//Use GoLang http functions so that we can get the content type header
	path := client.baseUrl + client.getResourcePath(ObjectsPath) + "/" + bucket + "/get/" + object
	headers := client.getRequestHeaders(accessToken)

	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, management.NewError(http.StatusBadGateway, "bad_gateway", err.Error())
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	httpClient := management.NewHttpClient()
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Error reading body: %s", err.Error()))
	}

	return &model.Data{
		Payload:     body,
		ContentType: response.Header.Get("Content-Type"),
	}, nil
}
