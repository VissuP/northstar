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

package provider

import (
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/portal/model"
)

// PortalProvider defines the basic interface that implementers of the portal api functionality must fulfill.
type PortalProvider interface {

	// CreateNotebook creates a new notebook
	CreateNotebook(token string, notebook *model.Notebook) (*model.Notebook, *management.Error)

	// UpdateNotebook updates an existing notebook.
	UpdateNotebook(token string, notebook *model.Notebook) *management.Error

	// ListNotebooks lists the existing notebooks for an account.
	ListNotebooks(token string) ([]model.Notebook, *management.Error)

	// GetNotebook retrieves the notebook matching a notebook ID.
	GetNotebook(token string, notebookID string) (*model.Notebook, *management.Error)

	// DeleteNotebook removes the notebook with the specified ID.
	DeleteNotebook(token string, notebookID string) *management.Error

	// ExecuteCell submits an execution request for a cell.
	ExecuteCell(token string, callbackURL string, cell *model.Cell) *management.Error

	// GetNotebookUsers returns a list of users that have access to the specified notebook
	GetNotebookUsers(token string, notebookId string) ([]model.User, *management.Error)

	// UpdateNotebookUsers updates the permissions of the users that have access to the specified notebook
	UpdateNotebookUsers(token string, notebookId string, users []model.User) *management.Error

	// ListPortfolios lists the existing portfolios for an account.
	ListPortfolios(token string) ([]model.Portfolio, *management.Error)

	// ListFiles lists the existing files for a portfolio.
	ListFiles(token string, portfolio string, prefix string, count int, marker string) ([]model.File, *management.Error)

	// GetFile gets the file for download.
	GetFile(token string, portfolio string, file string) (*model.Data, *management.Error)

	// CreateTransformation creates a new transformation.
	CreateTransformation(token string, transformation *model.Transformation) (*model.Transformation, *management.Error)

	// UpdateTransformation updates an existing Transformation.
	UpdateTransformation(token string, transformation *model.Transformation) *management.Error

	// ListTransformations lists the Transformations owned by an account.
	ListTransformations(token string) ([]model.Transformation, *management.Error)

	// GetTransformation retrieves the specified Transformation.
	GetTransformation(token string, transformationId string) (*model.Transformation, *management.Error)

	// GetTransformationResults returns the collection of execution results for the specified Transformation.
	GetTransformationResults(token string, transformationID string) ([]model.Output, *management.Error)

	// DeleteTransformation deletes the specified transformationId.
	DeleteTransformation(token string, transformationId string) *management.Error

	// ExecuteTransformation submits an execution request for a transformation.
	ExecuteTransformation(token string, callbackURL string, transformation *model.Transformation) *management.Error

	// Schedule
	CreateSchedule(token string, transformationId string, schedule *model.Schedule) *management.Error

	// GetSchedule retrieves the schedule for the specified transformationId.
	GetSchedule(token string, transformationId string) (*model.Schedule, *management.Error)

	// DeleteSchedule deletes the specified transformationId.
	DeleteSchedule(token string, transformationId string) *management.Error

	// GetSchemas retrieves the event schemas from thingspace via northstarapi
	GetScheduleEventSchemas(token string) ([]model.ScheduleEventSchema, *management.Error)

	// QueryUsers returns a list of Thingspace users filtered according to the provided query
	QueryUsers(token string, user *model.User) ([]model.User, *management.Error)

	// ProcessEvent parses the specified payload based on type.
	ProcessEvent(id, payloadType string, payload []byte) (*model.Event, *management.Error)

	// Templates
	ListTemplates(token string) ([]model.Template, *management.Error)
}
