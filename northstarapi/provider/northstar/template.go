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
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	cMap "github.com/orcaman/concurrent-map"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	templates "github.com/verizonlabs/northstar/data/templates/client"
	templatesModel "github.com/verizonlabs/northstar/data/templates/model"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Defines the type used to support operations on NorthStar notebooks.
type NorthStarTemplatesProvider struct {
	templatesClient *templates.TemplatesClient
	templateCache   cMap.ConcurrentMap
}

// Returns a new NorthStar notebook provider.
func NewNorthStarTemplatesProvider() (*NorthStarTemplatesProvider, error) {
	mlog.Info("NewNorthStarTemplatesProvider")

	// Create template client.
	templatesClient, err := templates.NewTemplatesClient()

	if err != nil {
		return nil, err
	}

	// Create the provider.
	provider := &NorthStarTemplatesProvider{
		templatesClient: templatesClient,
		templateCache:   cMap.New(),
	}

	return provider, nil
}

// Creates a new custom user template.
func (provider *NorthStarTemplatesProvider) Create(user *model.User, template *model.Template) (*model.Template, *management.Error) {
	mlog.Debug("Create")

	// Create external template.
	externalTemplate, err := provider.toExternalTemplate(user, template)
	if err != nil {
		mlog.Error("To external template returned error: %v", err)
		return nil, model.ErrorToExternalTemplate
	}

	// If template published settings not specified, set to private.
	if externalTemplate.Published == templatesModel.NotSet {
		externalTemplate.Published = templatesModel.Private
	}

	createTemplate, mErr := provider.templatesClient.CreateTemplate(externalTemplate)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Create template returned error: %v", mErr))
	}

	// Update template information before returning.
	template.Id = createTemplate.Id
	template.CreatedOn = createTemplate.CreatedOn.Format(time.RFC3339)

	return template, nil
}

// Updates a custom user template.
func (provider *NorthStarTemplatesProvider) Update(user *model.User, template *model.Template) *management.Error {
	mlog.Debug("Update")

	// Get external template.
	externalTemplate, mErr := provider.templatesClient.GetTemplate(template.Id)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Get template returned error: %v", mErr))
	}

	// If the user does not own the template, return error. Note that only the
	// owner can modify the template.
	if externalTemplate.UserId != user.Id {
		return management.ErrorForbidden
	}

	// Create external template.
	externalTemplate, err := provider.toExternalTemplate(user, template)

	if err != nil {
		mlog.Error("To external template returned error: %v", err)
		return model.ErrorToExternalTemplate
	}

	if _, mErr := provider.templatesClient.UpdateTemplate(externalTemplate.Id, externalTemplate); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Update template returned error: %v", mErr))
	}

	return nil
}

// Returns the custom user template with the specified id.
func (provider *NorthStarTemplatesProvider) Get(user *model.User, templateId string) (*model.Template, *management.Error) {
	mlog.Debug("Get")

	externalTemplate, mErr := provider.templatesClient.GetTemplate(templateId)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get template returned error: %v", mErr))
	}

	// If the template is not published, make sure the user owns it.
	if externalTemplate.Published != templatesModel.Published {
		// Validate user own the template.
		if externalTemplate.UserId != user.Id {
			return nil, management.ErrorForbidden
		}
	}

	template, err := provider.fromExternalTemplate(externalTemplate)

	if err != nil {
		return nil, model.ErrorFromExternalTemplate
	}

	return template, nil
}

//TemplateExists verifies that a template already exists
func (provider *NorthStarTemplatesProvider) TemplateExists(code string) *management.Error {
	mlog.Debug("VerifyTemplate")
	hash := provider.getHash(code)
	mlog.Debug("Looking for hash: %s", hash)
	mapTemplate, found := provider.templateCache.Get(hash)
	if found {
		mlog.Debug("Hash found in map: %s matches template: %s", hash, mapTemplate)
		return nil
	}

	query := &templatesModel.Query{
		Hash: hash,
	}
	retrievedTemplates, serviceErr := provider.templatesClient.QueryTemplate(query)
	if serviceErr != nil {
		mlog.Error("Failed to verify template with error: %s", serviceErr.Description)
		return management.GetInternalError("Code verification failed")
	}

	if len(retrievedTemplates) == 0 {
		mlog.Error("No templates returned for hash %s", hash)
		return management.GetInternalError("Code verification failed")
	}

	provider.templateCache.Set(hash, retrievedTemplates[0])

	return nil
}

//getHash generates a hash for a cell
func (provider *NorthStarTemplatesProvider) getHash(code string) string {
	mlog.Debug("getHash")

	hash := sha1.Sum([]byte(code))
	bytes := hash[:]
	encodedHash := base64.StdEncoding.EncodeToString(bytes)
	return encodedHash
}

// Returns the list of custom user templates associated with the authenticated user.
func (provider *NorthStarTemplatesProvider) List(user *model.User) ([]model.Template, *management.Error) {
	mlog.Debug("List")

	// Query user templates.
	query := &templatesModel.Query{
		UserId: user.Id,
	}

	externalTemplates, mErr := provider.templatesClient.QueryTemplate(query)

	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get (user) templates returned error: %v", mErr))
	}

	var templates []model.Template

	for _, externalTemplate := range externalTemplates {
		template, err := provider.fromExternalTemplate(&externalTemplate)

		if err != nil {
			mlog.Error("Failed to create external template: %v", template)
			continue
		}

		templates = append(templates, *template)
	}

	// Query public templates.
	query = &templatesModel.Query{
		Published: templatesModel.Published,
	}

	publicTemplates, mErr := provider.templatesClient.QueryTemplate(query)
	if mErr != nil {
		return nil, management.GetExternalError(fmt.Sprintf("Get (public) templates returned error: %v", mErr))
	}

	for _, publicTemplate := range publicTemplates {
		template, err := provider.fromExternalTemplate(&publicTemplate)

		if err != nil {
			mlog.Error("Failed to create public template: %v", template)
			continue
		}

		templates = append(templates, *template)
	}

	return templates, nil
}

// Deletes the template with the specified id.
func (provider *NorthStarTemplatesProvider) Delete(user *model.User, templateId string) *management.Error {
	mlog.Debug("Delete")

	externalTemplate, mErr := provider.templatesClient.GetTemplate(templateId)

	if mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Get template returned error: %v", mErr))
	}

	// Validate user own the template.
	if externalTemplate.UserId != user.Id {
		return management.ErrorForbidden
	}

	if mErr := provider.templatesClient.DeleteTemplate(templateId); mErr != nil {
		return management.GetExternalError(fmt.Sprintf("Delete template returned error: %v", mErr))
	}

	return nil
}

// Helper method used to translate portal api model to data service model.
func (provider *NorthStarTemplatesProvider) toExternalTemplate(user *model.User, template *model.Template) (*templatesModel.Template, error) {
	mlog.Debug("toExternalTemplate")

	// Note that we store template content as JSON.
	data, err := json.Marshal(template.Data)

	if err != nil {
		return nil, err
	}

	// Create the data service representation.
	externalTemplate := &templatesModel.Template{
		Id:          template.Id,
		UserId:      user.Id,
		AccountId:   user.AccountId,
		Version:     template.Etag,
		Name:        template.Name,
		Description: template.Description,
		Type:        template.Type,
		Data:        string(data),
		Published:   template.Published,
	}

	return externalTemplate, nil
}

// Helper method used to translate data service model to portal api model.
func (provider *NorthStarTemplatesProvider) fromExternalTemplate(externalTemplate *templatesModel.Template) (*model.Template, error) {
	mlog.Debug("fromExternalTemplate")

	template := &model.Template{
		Id:          externalTemplate.Id,
		Etag:        externalTemplate.Version,
		CreatedOn:   externalTemplate.CreatedOn.Format(time.RFC3339),
		Name:        externalTemplate.Name,
		Description: externalTemplate.Description,
		Type:        externalTemplate.Type,
		Published:   externalTemplate.Published,
	}

	// Unmarshal the data based on cell type.
	if template.Type == model.CellTemplateType {
		cell := &model.Cell{}

		if err := json.Unmarshal([]byte(externalTemplate.Data), &cell); err != nil {
			return nil, err
		}

		template.Data = cell
	} else if template.Type == model.NotebookTemplateType {
		notebook := &model.Notebook{}

		if err := json.Unmarshal([]byte(externalTemplate.Data), &notebook); err != nil {
			return nil, err
		}

		template.Data = notebook
	} else {
		mlog.Info("Warning, invalid or unsupported template type for id %s.", externalTemplate.Id)
		template.Type = model.UnknownTemplateType
	}

	return template, nil
}
