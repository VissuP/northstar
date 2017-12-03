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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/dpe-stream/config"
	"github.com/verizonlabs/northstar/pkg/marathon"
)

type MarathonCluster struct {
	marathonClient marathon.NorthstarMarathonClient
}

func NewMarathonCluster() (*MarathonCluster, error) {
	marathonClient, err := marathon.NewMarathonClient()
	if err != nil {
		return nil, err
	}

	return &MarathonCluster{marathonClient: marathonClient}, nil
}

func (m *MarathonCluster) StartJob(job *StartJob) error {
	mlog.Info("Starting job %s for account %s", job.JobId, job.AccountId)

	for i := 0; i < job.Instances; i++ {
		mlog.Debug("Starting worker %d", i)
		app, err := marathon.GetApplicationFromJson(config.WorkerMarathonJson)
		if err != nil {
			mlog.Error("Failed to get application from json: %v", err)
			return err
		}

		app.Name(getWorkerName(job.AccountId, job.JobId, i))
		app.Args = &[]string{"/usr/local/bin/dpe-stream worker"}
		env := *app.Env

		out, err := json.Marshal(job)
		if err != nil {
			mlog.Error("Failed to marshal job: %v", err)
			return err
		}

		env["DPE_STREAM_WORKER_JOB"] = b64.StdEncoding.EncodeToString(out)
		mlog.Debug("Stream job base64: %v", env["DPE_STREAM_WORKER_JOB"])

		err = m.marathonClient.CreateApplication(app)
		if err != nil {
			mlog.Error("Failed to create application: %v", err)
			return err
		}
	}

	return nil
}

func (m *MarathonCluster) StopJob(accountId string, jobId string) error {
	groupName := fmt.Sprintf("/%s/%s/dpe-stream-jobs/%s/%s",
		os.Getenv("MON_GROUP"), os.Getenv("ENV"), accountId, jobId)
	mlog.Debug("Group name: %v", groupName)
	return m.marathonClient.DeleteGroup(groupName)
}

func getWorkerName(accountId, jobId string, index int) string {
	appName := fmt.Sprintf("/%s/%s/dpe-stream-jobs/%s/%s/worker-%d",
		os.Getenv("MON_GROUP"), os.Getenv("ENV"), accountId, jobId, index)
	mlog.Debug("Worker name: %v", appName)
	return appName
}
