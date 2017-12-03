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
	"testing"
	"time"

	"github.com/satori/go.uuid"
	jobDataModel "github.com/verizonlabs/northstar/data/cron/model"
	"github.com/verizonlabs/northstar/processing/snippets/client"
)

const jobsCount = 10

func TestConvertJobData(t *testing.T) {
	jobsData := jobDataToCovert()
	jobs := ConvertJobDataArr(jobsData, client.SnippetsClientMock{})

	if len(jobsData) != len(jobs) {
		t.Errorf("Expected length of arrays to be the same, length of jobs data array: %v and length of jobs array: %v",
			len(jobsData), len(jobs))
	}
	for i, job := range jobs {
		if jobsData[i].Id != job.Id ||
			jobsData[i].AccountId != job.AccountId ||
			jobsData[i].Name != job.Name ||
			jobsData[i].SnippetId != job.SnippetId ||
			jobsData[i].Schedule != job.Schedule ||
			jobsData[i].Disabled != job.Disabled {
			t.Fail()
		}
	}
}

func jobDataToCovert() []*jobDataModel.JobData {
	var jobsData []*jobDataModel.JobData
	for i := 0; i < jobsCount; i++ {
		jobsData = append(jobsData, &jobDataModel.JobData{
			Id:        uuid.NewV4().String(),
			AccountId: uuid.NewV4().String(),
			Name:      uuid.NewV4().String(),
			SnippetId: uuid.NewV4().String(),
			Schedule:  uuid.NewV4().String(),
			Disabled:  false,
			UpdatedOn: time.Now(),
		})
	}

	return jobsData
}
