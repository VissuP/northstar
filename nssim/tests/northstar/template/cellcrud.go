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
	tests.Register(tests.TestId(config.TemplateCell), NewTemplateCellCRUDTest)
}

type TemplateCellCRUDTest struct {
	*northstar.NorthstarApiBaseTest
}

func NewTemplateCellCRUDTest() (tests.Test, error) {
	nsapiBase, err := northstar.NewNorthstarApiBaseTest()
	if err != nil {
		return nil, err
	}

	return &TemplateCellCRUDTest{
		NorthstarApiBaseTest: nsapiBase,
	}, nil
}

var (
	//steps
	listTemplatesStep              utils.Step = "List templates."
	createPublicCellTemplate       utils.Step = "Create public template"
	getPrivateCellUnauthorizedStep utils.Step = "Get private cell as unauthorized user."
	getPublicCellSecondUserStep    utils.Step = "Get public cell as second user."

	//common variables
	privateCellTemplateName = "Private Test Cell Template 1"
	publicCellTemplateName  = "Public Test Cell Template 1"
)

func (test *TemplateCellCRUDTest) Execute(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogDebug("Execute -- TemplateCellCRUDTest")

	//Get an initial count of templates
	initialTemplates, err := test.NorthstarApiClient.ListTemplates(test.Users[0].Token.AccessToken)
	if err != nil {
		return logs.LogStep(listTemplatesStep, false, "Error, could not list templates.")
	}

	initialTemplatesCount := len(initialTemplates)
	logs.LogStep(listTemplatesStep, true, "Retrieved (%d) initial templates: %+v", initialTemplatesCount, initialTemplates)

	//Test creating a private template
	privateCellTemplate := &northstarApiModel.Template{
		Name:        privateCellTemplateName,
		Description: "A cell template used for testing",
		Type:        northstarApiModel.CellTemplateType,
		Data: northstarApiModel.Cell{
			Name: "Number Card",
			Input: northstarApiModel.Input{
				Type:       "Code",
				Language:   "lua",
				EntryPoint: "main",
				Body:       `local output = require("nsOutput")`,
			},
		},
		Published: northstarApiModel.Private,
	}
	createdPrivateCellTemplate, err := test.NorthstarApiClient.CreateTemplate(test.Users[0].Token.AccessToken, privateCellTemplate)
	if err != nil {
		return logs.LogStep(createPrivateTemplateStep, false, "Error, could not create private cell template. Error was: %s", err)
	}
	logs.LogStep(createPrivateTemplateStep, true, "Successfully. Results: %+v", createdPrivateCellTemplate)
	defer test.DeleteTemplate(logs, test.Users[0].Token.AccessToken, createdPrivateCellTemplate.Id)

	//Test creating a public template
	publicCellTemplate := &northstarApiModel.Template{
		Name:        publicCellTemplateName,
		Description: "A cell template used for testing",
		Type:        northstarApiModel.CellTemplateType,
		Data: northstarApiModel.Cell{
			Name: "Number Card",
			Input: northstarApiModel.Input{
				Type:       "Code",
				Language:   "lua",
				EntryPoint: "main",
				Body:       `local output = require("nsOutput")`,
			},
		},
		Published: northstarApiModel.Published,
	}
	createdPublicCellTemplate, err := test.NorthstarApiClient.CreateTemplate(test.Users[0].Token.AccessToken, publicCellTemplate)
	if err != nil {
		return logs.LogStep(createPublicCellTemplate, false, "Error, could not create public cell template. Error was: %s", err)
	}
	logs.LogStep(createPublicCellTemplate, true, "Success. Template: %+v", createdPublicCellTemplate)
	defer test.DeleteTemplate(logs, test.Users[0].Token.AccessToken, createdPublicCellTemplate.Id)

	//Sleep to let cassandra populate
	time.Sleep(config.DEFAULT_CASSANDRA_SLEEP)

	// Check to make sure we can retrieve our templates
	_, err = test.NorthstarApiClient.GetTemplate(test.Users[0].Token.AccessToken, createdPrivateCellTemplate.Id)
	if err != nil {
		return logs.LogStep(getPrivateTemplateStep, false, "Couldn't get created private template. Error was: %s", err)
	}
	logs.LogStep(getPrivateTemplateStep, true, "Success.")

	_, err = test.NorthstarApiClient.GetTemplate(test.Users[0].Token.AccessToken, createdPublicCellTemplate.Id)
	if err != nil {
		return logs.LogStep(getPublicTemplateStep, false, "Couldn't get created public template. Error was: %s", err)
	}
	logs.LogStep(getPublicTemplateStep, true, "Success.")

	//Now check a non-owner user
	_, err = test.NorthstarApiClient.GetTemplate(test.Users[1].Token.AccessToken, createdPrivateCellTemplate.Id)
	if err == nil {
		return logs.LogStep(getPrivateCellUnauthorizedStep, false, "Error, unauthorized user was able to retrieve private cell: %s ", createdPrivateCellTemplate)
	}
	logs.LogStep(getPrivateCellUnauthorizedStep, true, "Prevented unauthorized access.")

	_, err = test.NorthstarApiClient.GetTemplate(test.Users[1].Token.AccessToken, createdPublicCellTemplate.Id)
	if err != nil {
		return logs.LogError("Error, secondary user was unable to retrieve public cell %s. Error was: %s", createdPublicCellTemplate, err)
	}
	logs.LogStep(getPublicCellSecondUserStep, true, "Success.")

	return nil
}

func (test *TemplateCellCRUDTest) CleanupHangingResource(logs *utils.Logger, callbacks *utils.ThreadSafeMap) error {
	logs.LogInfo("CleanupHangingResource")

	templates, err := test.NorthstarApiClient.ListTemplates(test.Users[0].Token.AccessToken)
	if err != nil {

	}

	for _, template := range templates {
		if template.Name == privateCellTemplateName || template.Name == publicCellTemplateName {
			if err := test.NorthstarApiClient.DeleteTemplate(test.Users[0].Token.AccessToken, template.Id); err != nil {
				logs.LogError("Failed to delete template with ID %s. %s", template.Id, err.Error())
			}
		}
	}
	return nil
}
