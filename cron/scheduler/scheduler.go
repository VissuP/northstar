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

type Scheduler interface {
	Start(jobs []*model.Job) error
	Stop()
	Restart(jobs []*model.Job) error
	GetEntry(job *model.Job) *cron.Entry
}

type JobScheduler struct {
	Cron    *cron.Cron
	Started bool
}

func NewScheduler() *JobScheduler {
	c := cron.New()
	return &JobScheduler{Cron: c, Started: false}
}

func (s *JobScheduler) Start(jobs []*model.Job) error {
	for _, job := range jobs {
		if job.Disabled {
			mlog.Info("Job %s is disabled", job.Id)
			continue
		}
		mlog.Info("Adding job %s, name: %s, disabled: %t, schedule: %s, snippet id: %s",
			job.Id, job.Name, job.Disabled, job.Schedule, job.SnippetId)
		err := s.Cron.AddJob(job.Schedule, job)
		if err != nil {
			return err
		}
	}
	s.Cron.Start()
	s.Started = true
	mlog.Info("Scheduled started")
	return nil
}

func (s *JobScheduler) Stop() {
	if s.Started {
		s.Cron.Stop()
		s.Started = false
		s.Cron = cron.New()
		mlog.Info("Scheduler stopped")
	}
}

func (s *JobScheduler) Restart(jobs []*model.Job) error {
	s.Stop()
	return s.Start(jobs)
}

func (s *JobScheduler) GetEntry(job *model.Job) *cron.Entry {
	for _, e := range s.Cron.Entries() {
		j, _ := e.Job.(*model.Job)
		if j.Name == job.Name {
			return e
		}
	}
	return nil
}
