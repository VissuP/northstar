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

package northstar

import (
	"fmt"

	"encoding/json"

	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/portal/model"
)

// ListTemplates returns templates associated with authenticated user.
func (provider *NorthStarPortalProvider) ListTemplates(token string) ([]model.Template, *management.Error) {
	mlog.Debug("ListTemplates")

	// Get the portal api (external) templates for user (i.e., from token).
	externalTemplates, mErr := provider.northstarApiClient.ListTemplates(token)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Failed to list templates with error: %v", mErr))
	}

	// Translate to portal templates.
	var templates []model.Template

	for _, externalTemplate := range externalTemplates {
		template, err := provider.FromExternalTemplate(&externalTemplate)

		// In case of errors, log the error but continue.
		if err != nil {
			mlog.Error("Failed to parse template with error: %+v", err)
			continue
		}

		templates = append(templates, *template)
	}

	return templates, nil
}

// FromExternalTemplate is a helper method used to translate portal api model to portal model.
func (provider *NorthStarPortalProvider) FromExternalTemplate(externalTemplate *northstarApiModel.Template) (*model.Template, error) {
	mlog.Debug("FromExternalTemplate: externalTemplate:%+v", externalTemplate)

	// Note that we need to marshal/unmarshal the data in order to generate
	// proper types.
	data, err := json.Marshal(externalTemplate.Data)

	if err != nil {
		return nil, fmt.Errorf("Marshal tempate data returned error: %v", err)
	}

	// Create template. Note that by default template type
	// gets set to unknown template.
	template := &model.Template{
		Id:        externalTemplate.Id,
		Name:      externalTemplate.Name,
		CreatedOn: externalTemplate.CreatedOn,
		Type:      model.UnknownTemplateType,
	}

	switch externalTemplate.Type {
	case northstarApiModel.CellTemplateType:
		mlog.Debug("Parsing cell template.")
		var externalCell northstarApiModel.Cell

		if err := json.Unmarshal(data, &externalCell); err != nil {
			return nil, fmt.Errorf("Unmarshal cell returned error: %v", err)
		}

		// Set template data to cell type.
		template.Type = model.CellTemplateType
		template.Data = fromExternalCell(&externalCell)
	case northstarApiModel.NotebookTemplateType:
		mlog.Debug("Parsing notebook template.")
		var externalNotebook northstarApiModel.Notebook

		if err := json.Unmarshal(data, &externalNotebook); err != nil {
			return nil, fmt.Errorf("Unmarshal notebook returned error: %v", err)
		}

		// Set template data to notebook type.
		template.Type = model.NotebookTemplateType
		template.Data = fromExternalNotebook(&externalNotebook)
	}

	return template, nil
}
