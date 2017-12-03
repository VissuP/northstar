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

package template

import (
	"time"

	northstarApiModel "github.com/verizonlabs/northstar/northstarapi/model"
	"github.com/verizonlabs/northstar/nssim/config"
	"github.com/verizonlabs/northstar/nssim/tests"
	"github.com/verizonlabs/northstar/nssim/tests/northstar"
	"github.com/verizonlabs/northstar/nssim/utils"
)

func init() {
	tests.Register(tests.TestId(config.TemplateNotebook), NewTemplateNotebookCRUDTest)
}

type TemplateNotebookCRUDTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewTemplateNotebookCRUDTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &TemplateNotebookCRUDTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//steps
	listTemplateStep                   utils.Step = "List templates"
	createPrivateTemplateStep          utils.Step = "Create private template."
	createPublicTemplateStep           utils.Step = "Create public template"
	getPrivateTemplateStep             utils.Step = "Get private template"
	getPublicTemplateStep              utils.Step = "Get public template"
	getPrivateTemplateUnauthorizedStep utils.Step = "Get private template as unauthorized user."
	getPublicTemplateSecondUser        utils.Step = "Get public template as second user."

	//common variables
	privateNotebookTemplateName = "Private Test Notebook Template 1"
	publicNotebookTemplateName  = "Public Test Notebook Template 1"
)

func (test *TemplateNotebookCRUDTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- TemplateNotebookCRUDTest")

	//Get an initial count of templates
	initialTemplates, err := test.NorthstarApiClient.ListTemplates(test.Users[0].Token.AccessToken)
	if err != nil {
		return logs.LogStep(listTemplateStep, false, "Error, could not list templates.")
	}

	initialTemplatesCount := len(initialTemplates)
	logs.LogStep(listTemplateStep, true, "Retrieved (%d) initial templates: %+v", initialTemplatesCount, initialTemplates)

	//Test creating a private template
	privateTemplate := &northstarApiModel.Template{
		Name:        privateNotebookTemplateName,
		Description: "A cell template used for testing",
		Type:        northstarApiModel.NotebookTemplateType,
		Data: northstarApiModel.Notebook{
			Name: "Notebook Template Test",
			Cells: []northstarApiModel.Cell{
				{
					Name: "Number Card",
					Input: northstarApiModel.Input{
						Type:       "Code",
						Language:   "lua",
						EntryPoint: "main",
						Body:       `local output = require("nsOutput")`,
					},
				},
			},
		},
		Published: northstarApiModel.Private,
	}
	createdPrivateTemplate, err := test.NorthstarApiClient.CreateTemplate(test.Users[0].Token.AccessToken, privateTemplate)
	if err != nil {
		return logs.LogStep(createPrivateTemplateStep, false, "Error, could not create private cell template. Error was: %s", err)
	}
	logs.LogStep(createPrivateTemplateStep, true, "Successfully created private cell template.")
	defer test.DeleteTemplate(logs, test.Users[0].Token.AccessToken, createdPrivateTemplate.Id)

	//Test creating a public template
	publicTemplate := &northstarApiModel.Template{
		Name:        publicNotebookTemplateName,
		Description: "A cell template used for testing",
		Type:        northstarApiModel.NotebookTemplateType,
		Data: northstarApiModel.Notebook{
			Name: "Notebook Template Test",
			Cells: []northstarApiModel.Cell{
				{
					Name: "Number Card",
					Input: northstarApiModel.Input{
						Type:       "Code",
						Language:   "lua",
						EntryPoint: "main",
						Body:       `local output = require("nsOutput")`,
					},
				},
			},
		},
		Published: northstarApiModel.Published,
	}

	createdPublicTemplate, err := test.NorthstarApiClient.CreateTemplate(test.Users[0].Token.AccessToken, publicTemplate)
	if err != nil {
		return logs.LogStep(createPublicTemplateStep, false, "Error, could not create public notebook template. Error was: %s", err)
	}
	logs.LogStep(createPublicTemplateStep, true, "Successfully created public cell template.")
	defer test.DeleteTemplate(logs, test.Users[0].Token.AccessToken, createdPublicTemplate.Id)

	time.Sleep(config.DEFAULT_CASSANDRA_SLEEP)

	// Check to make sure we can retrieve our templates
	_, err = test.NorthstarApiClient.GetTemplate(test.Users[0].Token.AccessToken, createdPrivateTemplate.Id)
	if err != nil {
		return logs.LogStep(getPrivateTemplateStep, false, "Couldn't get created private template. Error was: %s", err)
	}
	logs.LogStep(getPrivateTemplateStep, true, "Success.")

	_, err = test.NorthstarApiClient.GetTemplate(test.Users[0].Token.AccessToken, createdPublicTemplate.Id)
	if err != nil {
		return logs.LogStep(getPublicTemplateStep, false, "Couldn't get created public template. Error was: %s", err)
	}
	logs.LogStep(getPublicTemplateStep, true, "Success.")

	//Now check a non-owner user
	_, err = test.NorthstarApiClient.GetTemplate(test.Users[1].Token.AccessToken, createdPrivateTemplate.Id)
	if err == nil {
		return logs.LogStep(getPrivateTemplateUnauthorizedStep, false, "Error, unauthorized user was able to retrieve private notebook: %+v", createdPrivateTemplate)
	}
	logs.LogStep(getPrivateTemplateUnauthorizedStep, true, "Prevented unauthorized access.")

	_, err = test.NorthstarApiClient.GetTemplate(test.Users[1].Token.AccessToken, createdPublicTemplate.Id)
	if err != nil {
		return logs.LogStep(getPublicTemplateSecondUser, false, "Error, secondary user was unable to retrieve public notebook %s. Error was: %s", createdPublicTemplate, err)
	}
	logs.LogStep(getPublicTemplateSecondUser, true, "Success.")

	return nil
}

func (test *TemplateNotebookCRUDTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	templates, err := test.NorthstarApiClient.ListTemplates(test.Users[0].Token.AccessToken)
	if err != nil {

	}

	for _, template := range templates {
		if template.Name == privateNotebookTemplateName || template.Name == publicNotebookTemplateName {
			if err := test.NorthstarApiClient.DeleteTemplate(test.Users[0].Token.AccessToken, template.Id); err != nil {
				logs.LogError("Failed to delete template with ID %s. %s", template.Id, err.Error())
			}
		}
	}
	return nil
}
