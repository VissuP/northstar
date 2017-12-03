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
	"time"
)

type JobData struct {
	Id          string    `json:"id,omitempty"`
	AccountId   string    `json:"accountId,omitempty"`
	Name        string    `json:"name,omitempty"`
	SnippetId   string    `json:"snippetId,omitempty"`
	Schedule    string    `json:"schedule,omitempty"`
	Disabled    bool      `json:"disabled,omitempty"`
	UpdatedOn   time.Time `json:"updatedOn,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (job *JobData) ValidateOnAdd() error {
	if job.Id == "" {
		return fmt.Errorf("Id is empty")
	}

	if job.AccountId == "" {
		return fmt.Errorf("Account id is empty")
	}

	if job.SnippetId == "" {
		return fmt.Errorf("Snippet ID is empty")
	}

	if job.Schedule == "" {
		return fmt.Errorf("Schedule is empty")
	}
	return nil
}

func (job *JobData) ValidateOnUpdate() error {
	if job.UpdatedOn.IsZero() {
		return fmt.Errorf("Updated on is empty")
	}

	return nil
}

func (job *JobData) Print() string {
	return fmt.Sprintf("ID: %s, "+
		"Name: %s, "+
		"Disabled: %t, "+
		"Schedule: %s, "+
		"SnippetId: %s, "+
		"Description: %s\n",
		job.Id, job.Name, job.Disabled, job.Schedule, job.SnippetId, job.Description)
}
