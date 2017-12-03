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
	"github.com/verizonlabs/northstar/data/templates/model"
	"github.com/verizonlabs/northstar/data/util"
)

const TemplatesUri = util.DataBasePath + "/templates"

// Create a new template entry.
func (client *TemplatesClient) CreateTemplate(template *model.Template) (*model.Template, *management.Error) {
	mlog.Debug("CreateTemplate")

	path := TemplatesUri
	response, mErr := client.lbClient.PostJSON(path, template)

	if mErr != nil {
		return nil, mErr
	}

	createdTemplate := &model.Template{}

	if err := json.Unmarshal(response, createdTemplate); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return createdTemplate, nil
}

// Update an existing template entry.
func (client *TemplatesClient) UpdateTemplate(templateId string, template *model.Template) (*model.Template, *management.Error) {
	mlog.Debug("UpdateTemplate")

	path := TemplatesUri + "/" + templateId
	response, mErr := client.lbClient.PutJSON(path, template)

	if mErr != nil {
		return nil, mErr
	}

	updatedTemplate := &model.Template{}

	if err := json.Unmarshal(response, updatedTemplate); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return updatedTemplate, nil
}

// Returns template for the specified id.
func (client *TemplatesClient) GetTemplate(templateId string) (*model.Template, *management.Error) {
	mlog.Debug("GetTemplate")

	path := TemplatesUri + "/" + templateId
	response, mErr := client.lbClient.Get(path)

	if mErr != nil {
		return nil, mErr
	}

	template := &model.Template{}

	if err := json.Unmarshal(response, template); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return template, nil
}

// Query the collection of template with specific criteria.
func (client *TemplatesClient) QueryTemplate(query *model.Query) ([]model.Template, *management.Error) {
	mlog.Debug("QueryTemplate: query: +%v", query)

	path := TemplatesUri + "/actions/query"
	response, mErr := client.lbClient.PostJSON(path, query)

	if mErr != nil {
		return nil, mErr
	}

	templates := make([]model.Template, 0)

	if err := json.Unmarshal(response, &templates); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return templates, nil
}

// Deletes an template entry.
func (client *TemplatesClient) DeleteTemplate(templateId string) *management.Error {
	mlog.Debug("DeleteTemplate")

	path := TemplatesUri + "/" + templateId

	if mErr := client.lbClient.Delete(path); mErr != nil {
		return mErr
	}

	return nil
}
