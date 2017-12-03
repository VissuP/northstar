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
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/northstarapi/model"
)

func TestCreateTemplate(t *testing.T) {

	// Check if functional test has been enable.
	if os.Getenv("FUNCTIONAL_TESTING_ENABLED") == "" {
		t.Skip("Skipping test; $FUNCTIONAL_TESTING_ENABLED not set.")
	}

	Convey("Test Create Template", t, func() {
		provider, err := NewNorthStarTemplatesProvider()
		So(err, ShouldBeNil)

		// Create template.
		tableTemplate := createTableTemplate()
		tableTemplate.Published = model.Private

		user1 := &model.User{
			Id: "03f126c9-6e7f-4518-aea4-d49d74acce9d",
		}

		createTemplate, mErr := provider.Create(user1, tableTemplate)

		So(mErr, ShouldBeNil)
		So(createTemplate, ShouldNotBeNil)
		So(createTemplate.Name, ShouldEqual, tableTemplate.Name)
		So(createTemplate.Description, ShouldEqual, tableTemplate.Description)
		So(createTemplate.Published, ShouldEqual, tableTemplate.Published)

		// Get Template
		getTemplate, mErr := provider.Get(user1, createTemplate.Id)

		So(mErr, ShouldBeNil)
		So(getTemplate, ShouldNotBeNil)
		So(getTemplate.Name, ShouldEqual, tableTemplate.Name)
		So(getTemplate.Description, ShouldEqual, tableTemplate.Description)

		// Update the template name
		getTemplate.Name = "newName"
		getTemplate.Description = "new description"
		mErr = provider.Update(user1, getTemplate)

		So(mErr, ShouldBeNil)

		// Get Template to validate changes
		updatedTemplate, mErr := provider.Get(user1, tableTemplate.Id)

		So(mErr, ShouldBeNil)
		So(updatedTemplate, ShouldNotBeNil)
		So(updatedTemplate.Name, ShouldEqual, "newName")
		So(updatedTemplate.Description, ShouldEqual, "new description")

		// List Templates
		queryTemplates, mErr := provider.List(user1)

		So(queryTemplates, ShouldNotBeNil)
		So(mErr, ShouldBeNil)
		So(len(queryTemplates) > 0, ShouldEqual, true)

		// Delete the template.
		mErr = provider.Delete(user1, createTemplate.Id)

		So(mErr, ShouldBeNil)
	})

	Convey("Test Create Template - Public and non-public templates", t, func() {
		provider, err := NewNorthStarTemplatesProvider()
		So(err, ShouldBeNil)

		// Create template.
		tableTemplate := createTableTemplate()
		numberTemplate := createNumberTemplate()

		user1 := &model.User{
			Id: "03f126c9-6e7f-4518-aea4-d49d74acce9d",
		}
		user2 := &model.User{
			Id: "f293ca28-60e9-49e5-9d01-0490328b195d",
		}

		// Create table template for user 1. As private.
		tableTemplate.Published = model.Private
		createTableTemplate, mErr := provider.Create(user1, tableTemplate)
		So(mErr, ShouldBeNil)
		So(createTableTemplate, ShouldNotBeNil)

		// Create number template for user 2. As published.
		numberTemplate.Published = model.Published
		createNumberTemplate, mErr := provider.Create(user2, numberTemplate)
		So(mErr, ShouldBeNil)
		So(createNumberTemplate, ShouldNotBeNil)

		// Attempt to get template 1 with user 2. This should fail with forbidden error.
		getTableTemplate, mErr := provider.Get(user2, createTableTemplate.Id)
		So(mErr, ShouldNotBeNil)
		So(getTableTemplate, ShouldBeNil)

		// Attempt to get template 2 with user 1. This should succeed.
		getNumberTemplate, mErr := provider.Get(user1, createNumberTemplate.Id)
		So(mErr, ShouldBeNil)
		So(getNumberTemplate, ShouldNotBeNil)

		// Update number template with user 1. This should fail.
		mErr = provider.Update(user1, createNumberTemplate)
		So(mErr, ShouldNotBeNil)

		// Delete number template with user 1. This should fail.
		mErr = provider.Delete(user1, createNumberTemplate.Id)
		So(mErr, ShouldNotBeNil)

		// Cleanp.
		mErr = provider.Delete(user1, createTableTemplate.Id)
		So(mErr, ShouldBeNil)

		mErr = provider.Delete(user2, numberTemplate.Id)
		So(mErr, ShouldBeNil)
	})
}

func TestListTemplate(t *testing.T) {

	// Check if functional test has been enable.
	if os.Getenv("FUNCTIONAL_TESTING_ENABLED") == "" {
		t.Skip("Skipping test; $FUNCTIONAL_TESTING_ENABLED not set.")
	}

	Convey("Test List Template - No User", t, func() {
		provider, err := NewNorthStarTemplatesProvider()
		So(err, ShouldBeNil)

		noUser := &model.User{}

		// List Templates
		queryTemplates, mErr := provider.List(noUser)

		So(queryTemplates, ShouldNotBeNil)
		So(mErr, ShouldBeNil)
		So(len(queryTemplates) > 0, ShouldEqual, true)

		mlog.Debug(fmt.Sprintf("queryTemplates: +%v", queryTemplates))
	})

	Convey("Test List Template - User with no templates", t, func() {
		provider, err := NewNorthStarTemplatesProvider()
		So(err, ShouldBeNil)

		noUser := &model.User{
			Id: "5bfba265-ad81-4fe3-8419-d842d985c804",
		}

		// List Templates. Note that this should return all public templates.
		queryTemplates, mErr := provider.List(noUser)

		So(queryTemplates, ShouldNotBeNil)
		So(mErr, ShouldBeNil)
		So(len(queryTemplates) > 0, ShouldEqual, true)

		mlog.Debug(fmt.Sprintf("queryTemplates: +%v", queryTemplates))
	})
}

// Helper method used to create a table example template.
func createTableTemplate() *model.Template {
	luaCode := `
	--
	-- Copyright 2017 Verizon. All rights reserved.
	-- See provided LICENSE file for use of this source code.
	--
	-- NorthStar SDK Example : Output Table
	--
	local output = require("nsOutput")

	function main()
		-- Read column names from input arguments.
		local label1 = context.Args["label1"]
		local label2 = context.Args["label2"]

		local table = {
			columns = {label1, label2},
			rows = {
				{3, 60},{10, 65},{12, 70},
				{15, 75},{20, 80},{15, 85},
				{12, 90},{10, 95},{3, 100}},
		}

		-- Generate table results.
		return output.table(table)
	end
	`
	template := &model.Template{
		Type:        "cell",
		Name:        "Table Example",
		Description: "This template demonstrates usage of the NorthStar SDK table function.",
		Data: model.Cell{
			Name: "Table Example",
			Input: model.Input{
				Type:     model.CodeCellType,
				Language: "lua",
				Arguments: map[string]interface{}{
					"label1": "frequency",
					"label2": "percent",
				},
				EntryPoint: "main",
				Body:       luaCode,
				Timeout:    180,
			},
			Settings: model.Settings{
				ShowCode:   false,
				ShowOutput: true,
			},
		},
		Published: model.Published,
	}

	return template
}

// Helper method used to create a number example template.
func createNumberTemplate() *model.Template {
	luaCode := `
	--
	-- Copyright 2017 Verizon. All rights reserved.
	-- See provided LICENSE file for use of this source code.
	--
	-- NorthStar SDK Example : Output Number
	--
	local output = require("nsOutput")

	function main()
		-- Read semantics from input argument.
		local semantic = context.Args["semantic"]

  		local intTest = {
    		type = "number",
    		value = 78.08,
    		semantic = semantic
  		}

		-- Generate number results.
  		return output.number(intTest)
	end
	`
	template := &model.Template{
		Type:        "cell",
		Name:        "Number Example",
		Description: "This template demonstrates usage of the NorthStar SDK number function.",
		Data: model.Cell{
			Name: "Number Example",
			Input: model.Input{
				Type:     model.CodeCellType,
				Language: "lua",
				Arguments: map[string]interface{}{
					"semantic": "temperature",
				},
				EntryPoint: "main",
				Body:       luaCode,
				Timeout:    180,
			},
			Settings: model.Settings{
				ShowCode:   false,
				ShowOutput: true,
			},
		},
		Published: model.Published,
	}

	return template
}
