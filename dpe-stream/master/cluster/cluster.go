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

package cluster

import (
	"fmt"
	"github.com/verizonlabs/northstar/dpe-stream/master/model"
)

type StartJob struct {
	AccountId    string           `json:"accountId,omitempty"`
	JobId        string           `json:"jobId,omitempty"`
	InvocationId string           `json:"invocationId,omitempty"`
	Memory       uint64           `json:"memory,omitempty"`
	Instances    int              `json:"instances,omitempty"`
	Source       model.Source     `json:"source,omitempty"`
	Functions    []model.Function `json:"functions,omitempty"`
}

func (j *StartJob) Validate() error {
	if j.AccountId == "" {
		return fmt.Errorf("Account id is empty")
	}

	if j.JobId == "" {
		return fmt.Errorf("Job id is empty")
	}

	if j.InvocationId == "" {
		return fmt.Errorf("Invocation id is empty")
	}

	if j.Source.Name == "" {
		return fmt.Errorf("Source name is empty")
	}

	if len(j.Functions) < 1 {
		return fmt.Errorf("Number of functions less than one")
	}

	return nil
}

type Cluster interface {
	StartJob(job *StartJob) error
	StopJob(accountId, jobId string) error
}
