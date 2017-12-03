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
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/notebooks/model"
	"github.com/verizonlabs/northstar/data/util"
)

const AccessURI = util.DataBasePath + "/access"

// Create a new access entry.
func (client *NotebooksClient) CreateAccess(access *model.Access) (*model.Access, *management.Error) {
	mlog.Debug("CreateAccess")

	response, mErr := client.lbClient.PostJSON(AccessURI, access)

	if mErr != nil {
		return nil, mErr
	}

	createdAccess := &model.Access{}

	if err := json.Unmarshal(response, createdAccess); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return createdAccess, nil
}

// Update an existing access entry.
func (client *NotebooksClient) UpdateAccess(accessId string,
	access *model.Access) (*model.Access, *management.Error) {
	mlog.Debug("UpdateAccess")

	path := AccessURI + "/" + accessId
	response, mErr := client.lbClient.PutJSON(path, access)

	if mErr != nil {
		return nil, mErr
	}

	updatedAccess := &model.Access{}

	if err := json.Unmarshal(response, updatedAccess); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return updatedAccess, nil
}

// Query the collection of access with specific criteria.
func (client *NotebooksClient) QueryAccess(query *model.Query) ([]model.Access, *management.Error) {
	mlog.Debug("QueryAccess")

	path := AccessURI + "/actions/query"
	response, mErr := client.lbClient.PostJSON(path, query)

	if mErr != nil {
		return nil, mErr
	}

	var access []model.Access

	if err := json.Unmarshal(response, &access); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return access, nil
}

// Deletes an access entry.
func (client *NotebooksClient) DeleteAccess(accessId string) *management.Error {
	mlog.Debug("DeleteAccess")

	path := AccessURI + "/" + accessId

	if mErr := client.lbClient.Delete(path); mErr != nil {
		return mErr
	}

	return nil
}
