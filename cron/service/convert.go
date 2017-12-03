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

package service

import (
	"github.com/verizonlabs/northstar/cron/model"
	jobDataModel "github.com/verizonlabs/northstar/data/cron/model"
	"github.com/verizonlabs/northstar/processing/snippets/client"
)

func ConvertJobData(job *jobDataModel.JobData, processing client.Client) *model.Job {
	return &model.Job{Id: job.Id,
		AccountId:        job.AccountId,
		Name:             job.Name,
		Disabled:         job.Disabled,
		SnippetId:        job.SnippetId,
		Schedule:         job.Schedule,
		ProcessingClient: processing}
}

func ConvertJobDataArr(jobs []*jobDataModel.JobData,
	processing client.Client) []*model.Job {
	output := make([]*model.Job, len(jobs))

	for index, job := range jobs {
		output[index] = ConvertJobData(job, processing)
	}

	return output
}
