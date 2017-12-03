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

// Creates a new access entry.
func (service *NotebookService) createAccess(context *gin.Context) {
	mlog.Debug("createNotebook")

	// Bind the expected data resource object.
	var access model.Access

	if err := context.Bind(&access); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrInsertAccess.Incr()
		return
	}

	// Generate resource uuid
	uuid, err := gocql.RandomUUID()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to generate resource UUID with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertAccess.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrInsertAccess.Incr()
		return
	}

	// Populate access fields.
	access.Id = uuid.String()
	access.CreatedOn = time.Now()

	// Insert data into database.
	if _, err := database.Insert(Keyspace, AccessTable).
		Param("id", access.Id).
		Param("createdon", access.CreatedOn).
		Param("accountid", access.AccountId).
		Param("userid", access.UserId).
		Param("permissions", access.Permission).
		Param("notebookid", access.NotebookId).
		Exec(session); err != nil {
		mlog.Error("Database insert returned error: %s", err.Error())
		context.JSON(management.ErrorInternal.HttpStatus, management.ErrorInternal)
		ErrInsertAccess.Incr()
		return
	}

	InsertAccess.Incr()
	context.JSON(http.StatusOK, access)
}

// Update the specified access entry.
func (service *NotebookService) updateAccess(context *gin.Context) {
	mlog.Debug("updateAccess")

	// Get the access id.
	accessId := strings.TrimSpace(context.Params.ByName("accessId"))

	// Bind the expected data resource object.
	var access model.Access

	if err := context.Bind(&access); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrUpdateAccess.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrUpdateAccess.Incr()
		return
	}

	if _, err := database.Update(Keyspace, AccessTable).
		Param("accountid", access.AccountId).
		Param("userid", access.UserId).
		Param("permissions", access.Permission).
		Param("notebookid", access.NotebookId).
		Where("id", accessId).
		Exec(session); err != nil {
		mlog.Error("Database insert returned error: %s", err.Error())
		context.JSON(management.ErrorInternal.HttpStatus, management.ErrorInternal)
		ErrUpdateAccess.Incr()
		return
	}

	UpdateAccess.Incr()
	context.JSON(http.StatusOK, access)
}

// Returns access information associated with a search criteria.
func (service *NotebookService) queryAccess(context *gin.Context) {
	mlog.Debug("queryAccess")

	// Bind the expected data resource.
	var query model.Query

	if err := context.Bind(&query); err != nil {
		mlog.Error("Failed to decode request body: %v", err)
		ErrGetAccess.Incr()
		return
	}

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetAccess.Incr()
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

	// If valid notebook id, add to where clause.
	if query.NotebookId != "" {
		wheres = wheres.Bind("notebookid", query.NotebookId)
	}

	// Execute the query.
	results := make([]model.Access, 0)
	access := new(model.Access)

	// Note that expectation is that this is a small collection.
	builder := database.Select(Keyspace, AccessTable).
		Value("id", &access.Id).
		Value("createdon", &access.CreatedOn).
		Value("accountid", &access.AccountId).
		Value("userid", &access.UserId).
		Value("notebookid", &access.NotebookId).
		Value("permissions", &access.Permission).
		Wheres(wheres...).
		AllowFiltering()
	iter := builder.Iter(session)
	for builder.Next(iter) {
		results = append(results, *access)
	}

	// Close iterator.
	if err := iter.Close(); err != nil {
		errorMessage := fmt.Sprintf("Failed to close query iterator with error: %s.", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrGetAccess.Incr()
		return
	}

	GetAccess.Incr()
	context.JSON(http.StatusOK, results)
}

// Deletes the specified access for a notebook.
func (service *NotebookService) deleteAccess(context *gin.Context) {
	mlog.Debug("deleteAccess")

	// Get the access id.
	accessId := strings.TrimSpace(context.Params.ByName("accessId"))

	// Get database session.
	session, err := getSession()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get database session with error: %s", err.Error())
		mlog.Error(errorMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errorMessage))
		ErrDelAccess.Incr()
		return
	}

	// Remove from database.
	if _, err := database.Delete(Keyspace, AccessTable).
		Where("id", accessId).
		Exec(session); err != nil {
		errMessage := fmt.Sprintf("Failed to delete access with id %s with error: %+v", accessId, err)
		mlog.Error(errMessage)
		context.JSON(http.StatusInternalServerError, management.GetInternalError(errMessage))
		ErrDelAccess.Incr()
		return
	}

	DelAccess.Incr()
	context.String(http.StatusNoContent, http.StatusText(http.StatusNoContent))
}
