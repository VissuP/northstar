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

package scheduler

import (
	"github.com/robfig/cron"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/cron/model"
)

type SchedulerMock struct{}

func (s SchedulerMock) Start(jobs []*model.Job) error {
	mlog.Info("Starting job")
	return nil
}

func (s SchedulerMock) Stop() {
	mlog.Info("Stopping job")
}

func (s SchedulerMock) Restart(jobs []*model.Job) error {
	mlog.Info("Restarting jobs")
	return nil
}

func (s SchedulerMock) GetEntry(job *model.Job) *cron.Entry {
	mlog.Info("Getting entry")
	return &cron.Entry{}
}
