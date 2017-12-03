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

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

const (
	TemplatesPath = "templates"
)

// CreateTemplate creates a new template object.
func (client *Client) CreateTemplate(accessToken string, template *model.Template) (*model.Template, *management.Error) {
	mlog.Debug("CreateTemplate")

	// Create template.
	path := client.getResourcePath(TemplatesPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.PostJSONWithHeaders(path, template, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resources.
	createdTemplate := &model.Template{}

	if err := json.Unmarshal(response, createdTemplate); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return createdTemplate, nil
}

// ListTemplates returns all template associated with the specified access token.
func (client *Client) ListTemplates(accessToken string) ([]model.Template, *management.Error) {
	mlog.Debug("ListTemplates")

	// Update template.
	path := client.getResourcePath(TemplatesPath)
	headers := client.getRequestHeaders(accessToken)

	response, mErr := management.GetWithHeaders(client.baseUrl, path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resources.
	templates := make([]model.Template, 0)

	if err := json.Unmarshal(response, &templates); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return templates, nil
}

// GetTemplate gets the template with the specified id.
func (client *Client) GetTemplate(accessToken, templateId string) (*model.Template, *management.Error) {
	mlog.Debug("GetTemplate")

	// Get template.
	path := client.getResourcePath(TemplatesPath) + "/" + templateId
	headers := client.getRequestHeaders(accessToken)

	response, mErr := client.lbClient.GetWithHeaders(path, headers)

	// If error, return.
	if mErr != nil {
		return nil, mErr
	}

	// Otherwise, return the created resource.
	template := &model.Template{}

	if err := json.Unmarshal(response, template); err != nil {
		return nil, management.GetInternalError(fmt.Sprintf("Failed to umarshal body with error: %v", err))
	}

	return template, nil
}

//  DeleteTemplate deletes the template with the specified id.
func (client *Client) DeleteTemplate(accessToken, templateId string) *management.Error {
	mlog.Debug("DeleteTemplate")

	// Update template.
	path := client.getResourcePath(TemplatesPath) + "/" + templateId
	headers := client.getRequestHeaders(accessToken)

	return client.lbClient.DeleteWithHeaders(path, headers)
}
