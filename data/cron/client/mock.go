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

package client

import (
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/data/cron/model"
)

type CronClientMock struct{}

func (c CronClientMock) AddJob(accountId string, data *model.JobData) *management.Error {
	mlog.Info("Adding job for account ID: %s", accountId)
	return nil
}

func (c CronClientMock) DeleteJob(accountId string, jobId string) *management.Error {
	mlog.Info("Adding job with ID: %s for account ID: %s", jobId, accountId)
	return nil
}

func (c CronClientMock) GetJob(accountId string,
	jobId string) (*model.JobData, *management.Error) {
	mlog.Info("Retrieving job with ID: %s and account ID: %s", jobId, accountId)
	return &model.JobData{}, nil
}

func (c CronClientMock) GetAllJobs() ([]*model.JobData, *management.Error) {
	mlog.Info("Retrieving all jobs")
	return []*model.JobData{}, nil
}

func (c CronClientMock) GetJobsByAccountId(accountId string) ([]*model.JobData, *management.Error) {
	mlog.Info("Retrieving all jobs under account ID: ", accountId)
	return []*model.JobData{}, nil
}

func (c CronClientMock) UpdateJob(accountId string,
	jobId string,
	update *model.JobData) *management.Error {
	mlog.Info("Updating job with ID %s and account ID: %s", jobId, accountId)
	return nil
}
