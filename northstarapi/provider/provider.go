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
	"github.com/verizonlabs/northstar/northstarapi/model"
)

// Defines the interface used to support account resource operations.
type AccountProvider interface {
	// Define common resource operations.
	GetAccountIdForLoginname(loginname string) (string, *management.Error)
	GetUser(loginname string) (*model.User, *management.Error)
	GetUserById(id string) (*model.User, *management.Error)
	SearchUsers(user *model.User) ([]model.User, *management.Error)
}

// Defines the interface used to support transformation
// resource operations.
type TransformationProvider interface {
	// Define common resource operations.
	Create(accountId string, transformation *model.Transformation) (*model.Transformation, *management.Error)
	Update(accountId string, transformation *model.Transformation) *management.Error
	Get(accountId string, transformationId string) (*model.Transformation, *management.Error)
	List(accountId string) ([]model.Transformation, *management.Error)
	Delete(accountId string, transformationId string) *management.Error

	// Defines the transformation schedule operations.
	CreateSchedule(accountId string, transformationId string, schedule *model.Schedule) *management.Error
	GetSchedule(accountId string, transformationId string) (*model.Schedule, *management.Error)
	DeleteSchedule(accountId string, transformationId string) *management.Error

	// Defines the transformation execution results.
	Results(accountID string, transformationID string) ([]model.Output, *management.Error)
}

// Defines the interface used to support notebook resource operations.
type NotebookProvider interface {
	//Create creates a new notebook
	Create(user *model.User, notebook *model.Notebook) (*model.Notebook, *management.Error)

	//Update updates the specified notebook
	Update(user *model.User, notebook *model.Notebook) *management.Error

	//Get retrieves the specified notebook
	Get(user *model.User, notebookId string) (*model.Notebook, *management.Error)

	//List retrieves the list of notebooks belonging to the user
	List(user *model.User) ([]model.Notebook, *management.Error)

	//Delete deletes the specified notebook
	Delete(user *model.User, notebookId string) *management.Error

	// GetUsers retreives the list of users and permissions for the specified notebook
	GetUsers(user *model.User, notebookId string) ([]model.User, *management.Error)

	//UpdateUsers updates the list of users and permissions for the specified notebook.
	UpdateUsers(user *model.User, notebookId string, users []model.User) *management.Error

	//GetExecutionInformation retrieves the execution information for the notebook
	GetExecutionInformation(user *model.User, notebookId string) (*model.User, *management.Error)
}

// Defines the interface used to support code executions.
type ExecutionProvider interface {
	//ExecuteCell executes the provided cell.
	ExecuteCell(user *model.User, cell *model.Cell, callback string) (*model.Cell, *management.Error)

	//ExecuteTransformation executes the provided transformation
	ExecuteTransformation(user *model.User, transformation *model.Transformation, callback string) (*model.Transformation, *management.Error)

	//ExecutionCallback processes the execution response
	ExecutionCallback(response *model.ExecutionResponse)

	//Execute triggers an execution
	Execute(accountId string, execution *model.ExecutionRequest) (*model.ExecutionRequest, *management.Error)

	//ListExecutions returns the executions related to that account
	ListExecutions(accountId string, limit int) ([]model.Output, *management.Error)

	//GetExecution returns the execution associated with the provided execution ID
	GetExecution(accountId string, executionId string) (*model.Output, *management.Error)

	//StopExecution stops the specified execution
	StopExecution(accountId string, executionId string) *management.Error
}

// Defines the interface used to support template operations.
type TemplateProvider interface {
	//Create creates a new template using the provided template.
	Create(user *model.User, template *model.Template) (*model.Template, *management.Error)

	//Update updates the specified template
	Update(user *model.User, template *model.Template) *management.Error

	//Get retrieves the specified template.
	Get(user *model.User, templateId string) (*model.Template, *management.Error)

	//List lists templates that the user has access to.
	List(user *model.User) ([]model.Template, *management.Error)

	//Delete deletes the specified template
	Delete(user *model.User, templateId string) *management.Error

	//TemplateExists verifies that a template already exists for a the supplied code
	TemplateExists(code string) *management.Error
}

//ObjectProvider defines the interface used to support bucket operations
type ObjectProvider interface {

	//GetBuckets returns the list of buckets belonging to the user.
	GetBuckets(user *model.User) ([]model.Bucket, *management.Error)

	//GetObjects returns the list of object in the bucket
	ListObjects(user *model.User, bucket string, prefix string) ([]model.Object, *management.Error)

	//GetObject returns the specified object
	GetObject(user *model.User, bucket string, path string) (*model.Data, *management.Error)
}

//StreamProvider defines the interface for supporting long running jobs
type StreamProvider interface {
	//ListStreams lists the jobs belonging to the user
	ListStreams(accountId string) ([]model.Stream, *management.Error)

	//GetStream retrieves the specified job
	GetStream(accountId string, jobId string) (*model.Stream, *management.Error)

	//RemoveStream removes the specified job
	RemoveStream(accountId string, jobId string) *management.Error
}
