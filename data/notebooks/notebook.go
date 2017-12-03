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

package notebooks

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
	"github.com/verizonlabs/northstar/data/notebooks/model"
)

// Creates a new notebook.
func (service *NotebookService) createNotebook(context *gin.Context) {
	mlog.Debug("createNotebook")

	// Bind the expected data resource object.
	var notebook model.Notebook

	if err := context.Bind(&notebook); err != nil {
		mlog.Error("Failed to decode request body with error: %v", err)
		ErrInsertNotebook.Incr()
		return
	}

	// Generate resource uuid
	uuid, err := gocql.RandomUUID()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to generate resource UUID with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertNotebook.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get DB session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertNotebook.Incr()
		return
	}

	// Populate notebook fields.
	notebook.Id = uuid.String()
	notebook.CreatedOn = time.Now()
	notebook.Version = gocql.TimeUUID().String()

	// Insert data into database.
	if _, err := database.Insert(Keyspace, NotebookTable).
		Param("id", notebook.Id).
		Param("createdon", notebook.CreatedOn).
		Param("version", notebook.Version).
		Param("data", notebook.Data).
		Exec(session); err != nil {
		mlog.Error("Failed to create notebook with error: %s", err.Error())
		context.JSON(management.ErrorInternal.HttpStatus, management.ErrorInternal)
		ErrInsertNotebook.Incr()
		return
	}

	InsertNotebook.Incr()
	context.JSON(http.StatusOK, notebook)
}

// Update the specified notebook.
func (service *NotebookService) updateNotebook(context *gin.Context) {
	mlog.Debug("updateNotebook")

	// Get the notebook id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	// Bind the expected data resource object.
	var notebook model.Notebook

	if err := context.Bind(&notebook); err != nil {
		mlog.Error("Failed to decode request body with error: %v", err)
		ErrUpdateNotebook.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrUpdateNotebook.Incr()
		return
	}

	// Update the version.
	currentVersion := notebook.Version
	notebook.Version = gocql.TimeUUID().String()

	// Create query builder. Note that we do not allow modification
	// of account id or owner id. If provided fields will be ignore.
	builder := database.Update(Keyspace, NotebookTable).
		Param("version", notebook.Version).
		Param("data", notebook.Data).
		Where("id", notebookId)

	// Note that version can be specified to support conditional updates.
	if currentVersion != "" {
		builder = builder.If("version", currentVersion)
	}

	// Execute the query.
	if updated, err := builder.Exec(session); err != nil {
		mlog.Error("Failed to update notebook with error: %s", err.Error())
		context.JSON(management.ErrorInternal.HttpStatus, management.ErrorInternal)
		ErrUpdateNotebook.Incr()
		return
	} else if updated == false {
		mlog.Error("Failed to update notebook with id %s with error: Not Found.", notebookId)
		context.JSON(management.ErrorNotFound.HttpStatus, management.ErrorNotFound)
		ErrUpdateNotebook.Incr()
		return
	}

	UpdateNotebook.Incr()
	context.JSON(http.StatusOK, notebook)
}

// Returns notebook with specified id.
func (service *NotebookService) getNotebook(context *gin.Context) {
	mlog.Debug("getNotebook")

	// Get the notebook id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetNotebook.Incr()
		return
	}

	// Get resource from database.
	var notebook model.Notebook

	if err := database.Select(Keyspace, NotebookTable).
		Value("id", &notebook.Id).
		Value("createdon", &notebook.CreatedOn).
		Value("version", &notebook.Version).
		Value("data", &notebook.Data).
		Where("id", notebookId).
		AllowFiltering().
		Scan(session); err != nil {
		mlog.Error("Failed to get notebook with id %s with error: Not Found.", notebookId)
		context.JSON(management.ErrorNotFound.HttpStatus, management.ErrorNotFound)
		ErrGetNotebook.Incr()
		return
	}

	GetNotebook.Incr()
	context.JSON(http.StatusOK, notebook)
}

// Deletes the notebook with the specified id.
func (service *NotebookService) deleteNotebook(context *gin.Context) {
	mlog.Debug("delete")

	// Get the notebook id.
	notebookId := strings.TrimSpace(context.Params.ByName("notebookId"))

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelNotebook.Incr()
		return
	}

	// Remove from database.
	success := true
	if success, err = database.Delete(Keyspace, NotebookTable).
		Where("id", notebookId).
		Exec(session); err != nil {
		errMessage := fmt.Sprintf("Failed to delete notebook with id %s with error: %+v", notebookId, err)
		mlog.Error(errMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errMessage))
		ErrDelNotebook.Incr()
		return
	}

	if !success {
		ErrDelNotebook.Incr()
		mlog.Error("Notebook %s not found.", notebookId)
		context.JSON(http.StatusInternalServerError, management.GetInternalError("Notebook not found."))
		return
	}

	DelNotebook.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
