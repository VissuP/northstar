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

package templates

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/verizonlabs/northstar/pkg/database"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/templates/model"
)

// Creates a new template entry.
func (service *TemplateService) createTemplate(context *gin.Context) {
	mlog.Debug("createTemplate")

	// Bind the expected data resource object.
	var template model.Template

	if err := context.Bind(&template); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		UserErrCreateTemplate.Incr()
		return
	}

	// Generate resource uuid
	uuid, err := gocql.RandomUUID()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to generate resource UUID with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		SvcErrCreateTemplate.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		SvcErrCreateTemplate.Incr()
		return
	}

	// Populate template fields.
	template.Id = uuid.String()
	template.CreatedOn = time.Now()
	template.Version = gocql.TimeUUID().String()

	// Insert data into database.
	if _, err := database.Insert(Keyspace, TemplateTable).
		Param("id", template.Id).
		Param("createdon", template.CreatedOn).
		Param("version", template.Version).
		Param("accountid", template.AccountId).
		Param("userid", template.UserId).
		Param("name", template.Name).
		Param("description", template.Description).
		Param("type", template.Type).
		Param("data", template.Data).
		Param("published", template.Published).
		Param("hash", template.Hash).
		Exec(session); err != nil {
		mlog.Error("Database insert returned error: %s", err.Error())
		context.JSON(http.StatusBadGateway, management.GetExternalError(fmt.Sprintf("Database insert returned error: %s", err.Error())))
		ExtErrCreateTemplate.Incr()
		return
	}

	CreateTemplate.Incr()
	context.JSON(http.StatusOK, template)
}

// Update the specified template entry.
func (service *TemplateService) updateTemplate(context *gin.Context) {
	mlog.Debug("updateTemplate")

	// Get the template id.
	templateId := strings.TrimSpace(context.Params.ByName("templateId"))

	// Bind the expected data resource object.
	var template model.Template

	if err := context.Bind(&template); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		SvcErrUpdateTemplate.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		SvcErrUpdateTemplate.Incr()
		return
	}

	// Create the builder
	builder := database.Update(Keyspace, TemplateTable).
		Param("name", template.Name).
		Param("description", template.Description).
		Param("type", template.Type).
		Param("data", template.Data).
		Param("hash", template.Hash)

	// Set the publish flag, if and only if provided.
	if template.Published != model.NotSet {
		builder = builder.Param("published", template.Published)
	}

	builder = builder.Where("id", templateId)

	// Update the template. Note that we only allow to update type and data.
	if _, err := builder.Exec(session); err != nil {
		mlog.Error("Database update returned error: %s", err.Error())
		context.JSON(http.StatusBadGateway, management.GetExternalError(fmt.Sprintf("Database update returned error: %s", err.Error())))
		ExtErrUpdateTemplate.Incr()
		return
	}

	UpdateTemplate.Incr()
	context.JSON(http.StatusOK, template)
}

// Returns template with specified id.
func (service *TemplateService) getTemplate(context *gin.Context) {
	mlog.Debug("getTemplate")

	// Get the template id.
	templateId := strings.TrimSpace(context.Params.ByName("templateId"))

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		SvcErrGetTemplate.Incr()
		return
	}

	// Get resource from database.
	var template model.Template

	if err := database.Select(Keyspace, TemplateTable).
		Value("id", &template.Id).
		Value("createdon", &template.CreatedOn).
		Value("version", &template.Version).
		Value("accountid", &template.AccountId).
		Value("userid", &template.UserId).
		Value("type", &template.Type).
		Value("name", &template.Name).
		Value("description", &template.Description).
		Value("data", &template.Data).
		Value("published", &template.Published).
		Value("hash", &template.Hash).
		Where("id", templateId).
		AllowFiltering().
		Scan(session); err != nil {
		mlog.Error("Database select returned error: %s.", err.Error())
		context.JSON(http.StatusBadGateway, management.GetExternalError(fmt.Sprintf("Database select returned error: %s", err.Error())))
		ExtErrGetTemplate.Incr()
		return
	}

	GetTemplate.Incr()
	context.JSON(http.StatusOK, template)
}

// Returns template information associated with a search criteria.
func (service *TemplateService) queryTemplate(context *gin.Context) {
	mlog.Debug("queryTemplate")

	// Bind the expected data resource.
	var query model.Query

	if err := context.Bind(&query); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		UserErrQueryTemplate.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		SvcErrQueryTemplate.Incr()
		return
	}

	// Build the where clause.
	var wheres database.Bindings

	// If valid account id, add to where clause.
	if query.AccountId != "" {
		wheres = wheres.Bind("accountid", query.AccountId)
	}

	// If valid user id, add to where clause.
	if query.UserId != "" {
		wheres = wheres.Bind("userid", query.UserId)
	}

	// If valid template id, add to where clause.
	if query.Type != "" {
		wheres = wheres.Bind("type", query.Type)
	}

	// If published, add to where clause.
	if query.Published != model.NotSet {
		wheres = wheres.Bind("published", query.Published)
	}

	//If hash, add to where clause
	if query.Hash != "" {
		wheres = wheres.Bind("hash", query.Hash)
	}

	// Execute the query.
	results := make([]model.Template, 0)
	template := new(model.Template)

	// Note that expectation is that this is a small collection.
	builder := database.Select(Keyspace, TemplateTable).
		Value("id", &template.Id).
		Value("createdon", &template.CreatedOn).
		Value("version", &template.Version).
		Value("accountid", &template.AccountId).
		Value("userid", &template.UserId).
		Value("name", &template.Name).
		Value("description", &template.Description).
		Value("type", &template.Type).
		Value("data", &template.Data).
		Value("published", &template.Published).
		Value("hash", &template.Hash).
		Wheres(wheres...).
		AllowFiltering()
	iter := builder.Iter(session)
	for builder.Next(iter) {
		results = append(results, *template)
	}

	// Close iterator.
	if err := iter.Close(); err != nil {
		mlog.Error("Database close returned error: %s", err.Error())
		context.JSON(http.StatusBadGateway, management.GetExternalError(fmt.Sprintf("Database close returned error: %s", err.Error())))
		ExtErrQueryTemplate.Incr()
		return
	}

	QueryTemplate.Incr()
	context.JSON(http.StatusOK, results)
}

// Deletes the specified template for a template.
func (service *TemplateService) deleteTemplate(context *gin.Context) {
	mlog.Debug("deleteTemplate")

	// Get the template id.
	templateId := strings.TrimSpace(context.Params.ByName("templateId"))

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		SvcErrDeleteTemplate.Incr()
		return
	}

	// Remove from database.
	success := true
	if success, err = database.Delete(Keyspace, TemplateTable).
		Where("id", templateId).
		Exec(session); err != nil {
		mlog.Error("Database delete returned error: %s", err.Error())
		context.JSON(http.StatusBadGateway, management.GetExternalError(fmt.Sprintf("Database delete returned error: %s", err.Error())))
		ExtErrDeleteTemplate.Incr()
		return
	}

	if !success {
		mlog.Error("Failed to delete template. Template %s not found.", templateId)
		SvcErrDeleteTemplate.Incr()
		context.JSON(http.StatusInternalServerError, management.GetInternalError("Template not found."))
		return
	}

	DelTemplate.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
