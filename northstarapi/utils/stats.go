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

package utils

import "github.com/verizonlabs/northstar/pkg/stats"

var (
	Stats             = stats.New("northstarapi")
	CreateTemplate    = Stats.NewCounter("CreateTemplate")
	ErrCreateTemplate = Stats.NewCounter("ErrCreateTemplate")
	ListTemplates     = Stats.NewCounter("ListTemplates")
	ErrListTemplates  = Stats.NewCounter("ErrListTemplates")
	GetTemplate       = Stats.NewCounter("GetTemplate")
	ErrGetTemplate    = Stats.NewCounter("ErrGetTemplate")
	DeleteTemplate    = Stats.NewCounter("DeleteTemplate")
	ErrDeleteTemplate = Stats.NewCounter("ErrDeleteTemplate")

	CreateNotebook         = Stats.NewCounter("CreateNotebook")
	ErrCreateNotebook      = Stats.NewCounter("ErrCreateNotebook")
	ListNotebooks          = Stats.NewCounter("ListNotebooks")
	ErrListNotebooks       = Stats.NewCounter("ErrListNotebooks")
	GetNotebook            = Stats.NewCounter("GetNotebook")
	ErrGetNotebook         = Stats.NewCounter("ErrGetNotebook")
	UpdateNotebook         = Stats.NewCounter("UpdateNotebook")
	ErrUpdateNotebook      = Stats.NewCounter("ErrUpdateNotebook")
	ExecuteNotebookCell    = Stats.NewCounter("ExecuteNotebookCell")
	ErrExecuteNotebookCell = Stats.NewCounter("ErrExecuteNotebookCell")
	ExecutionCallback      = Stats.NewCounter("ExecutionCallback")
	ErrExecutionCallback   = Stats.NewCounter("ErrExecutionCallback")
	DeleteNotebook         = Stats.NewCounter("DeleteNotebook")
	ErrDeleteNotebook      = Stats.NewCounter("ErrDeleteNotebook")
	GetUsers               = Stats.NewCounter("GetUsers")
	ErrGetUsers            = Stats.NewCounter("ErrGetUsers")
	UpdateUsers            = Stats.NewCounter("UpdateUsers")
	ErrUpdateUsers         = Stats.NewCounter("ErrUpdateUsers")
	ExecuteCell            = Stats.NewCounter("ExecuteCell")
	ErrExecuteCell         = Stats.NewCounter("ErrExecuteCell")

	TriggerExecution    = Stats.NewCounter("TriggerExecution")
	ErrTriggerExecution = Stats.NewCounter("ErrTriggerExecution")
	GetExecution        = Stats.NewCounter("GetExecution")
	ErrGetExecution     = Stats.NewCounter("ErrGetExecution")
	ListExecutions      = Stats.NewCounter("ListExecutions")
	ErrListExecutions   = Stats.NewCounter("ErrListExecutions")
	StopExecution       = Stats.NewCounter("StopExecution")
	ErrStopExecution    = Stats.NewCounter("ErrStopExecution")

	CreateTransformation     = Stats.NewCounter("CreateTransformation")
	ErrCreateTransformation  = Stats.NewCounter("ErrCreateTransformation")
	ListTransformations      = Stats.NewCounter("ListTransformations")
	ErrListTransformations   = Stats.NewCounter("ErrListTransformations")
	UpdateTransformation     = Stats.NewCounter("UpdateTransformation")
	ErrUpdateTransformation  = Stats.NewCounter("ErrUpdateTransformation")
	GetTransformation        = Stats.NewCounter("GetTransformation")
	ErrGetTransformation     = Stats.NewCounter("ErrGetTransformation")
	TransformationResults    = Stats.NewCounter("TransformationResults")
	ErrTransformationResults = Stats.NewCounter("ErrTransformationResults")
	DeleteTransformation     = Stats.NewCounter("DeleteTransformation")
	ErrDeleteTransformation  = Stats.NewCounter("ErrDeleteTransformation")
	ExecuteTransformation    = Stats.NewCounter("ExecuteTransformation")
	ErrExecuteTransformation = Stats.NewCounter("ErrExecuteTransformation")

	CreateSchedule    = Stats.NewCounter("CreateSchedule")
	ErrCreateSchedule = Stats.NewCounter("ErrCreateSchedule")
	GetSchedule       = Stats.NewCounter("GetSchedule")
	ErrGetSchedule    = Stats.NewCounter("ErrGetSchedule")
	DeleteSchedule    = Stats.NewCounter("DeleteSchedule")
	ErrDeleteSchedule = Stats.NewCounter("ErrDeleteSchedule")

	SearchUsers    = Stats.NewCounter("SearchUsers")
	ErrSearchUsers = Stats.NewCounter("ErrSearchUsers")

	ListBuckets    = Stats.NewCounter("ListBuckets")
	ErrListBuckets = Stats.NewCounter("ErrListBuckets")
	ListObjects    = Stats.NewCounter("ListObjects")
	ErrListObjects = Stats.NewCounter("ErrListObjects")
	GetObject      = Stats.NewCounter("GetObject")
	ErrGetObject   = Stats.NewCounter("ErrGetObject")

	ListStreams     = Stats.NewCounter("ListStreams")
	ErrListStreams  = Stats.NewCounter("ErrListStreams")
	GetStream       = Stats.NewCounter("GetStream")
	ErrGetStream    = Stats.NewCounter("ErrGetStream")
	RemoveStream    = Stats.NewCounter("RemoveStream")
	ErrRemoveStream = Stats.NewCounter("ErrRemoveStream")
)
