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

package model

import (
	"fmt"

	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/processing/snippets/client"
	"github.com/verizonlabs/northstar/processing/snippets/model"
)

const (
	MissingName      = "Name is empty"
	MissingSchedule  = "Schedule is empty"
	MissingSnippetId = "Snippet id is empty"
)

type Job struct {
	AccountId        string        `json:"-"`
	Id               string        `json:"jobId,omitempty"`
	Name             string        `json:"name,omitempty"`
	Disabled         bool          `json:"disabled,omitempty"`
	Schedule         string        `json:"schedule,omitempty"`
	SnippetId        string        `json:"snippetId,omitempty"`
	Description      string        `json:"description,omitempty"`
	ProcessingClient client.Client `json:"-"`
}

func (job Job) Validate() error {
	if job.Name == "" {
		return fmt.Errorf(MissingName)
	}

	if job.Schedule == "" {
		return fmt.Errorf(MissingSchedule)
	}

	if job.SnippetId == "" {
		return fmt.Errorf(MissingSnippetId)
	}

	return nil
}

func (job *Job) Run() {
	if job.Disabled {
		mlog.Info("Job %s is disabled", job.Id)
		return
	}

	if job.ProcessingClient == nil {
		mlog.Error("Processing client is null")
		return
	}

	mlog.Info("Running job %s, name: %s, schedule: %s, snippet id: %s",
		job.Id, job.Name, job.Schedule, job.SnippetId)

	snippet := model.Snippet{SnippetId: job.SnippetId, Options: model.Options{}}
	id, err := job.ProcessingClient.StartSnippet(job.AccountId, &snippet)
	if err != nil {
		mlog.Error("Failed to invoke snippet by id: %v", err)
		return
	}

	mlog.Info("Snippet %s from job %s invoked with invocation id %s", job.SnippetId, job.Id, id)
}
