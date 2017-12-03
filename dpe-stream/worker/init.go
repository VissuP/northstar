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

package worker

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/service_master"
	"github.com/verizonlabs/northstar/dpe-stream/config"
	"github.com/verizonlabs/northstar/dpe-stream/master/cluster"
	"github.com/verizonlabs/northstar/dpe-stream/master/connection"
	"github.com/verizonlabs/northstar/dpe-stream/master/model"
	"github.com/verizonlabs/northstar/dpe-stream/worker/events"
	"github.com/verizonlabs/northstar/dpe-stream/worker/source/kafka"
	"github.com/verizonlabs/northstar/dpe-stream/worker/stats"
)

func StartWorker() error {
	svcMaster := service_master.New(config.NumThreads, config.WorkerQueueCapacity)
	eventsProducer, err := events.NewKafkaEventsProducer()
	if err != nil {
		return err
	}

	job, err := getStreamingJob()
	if err != nil {
		stats.ErrGetJob.Incr()
		mlog.Error("Failed to get job description: %v", job)
		return err
	}

	err = job.Validate()
	if err != nil {
		stats.ErrValidateJob.Incr()
		return err
	}

	switch job.Source.Name {
	case model.SOURCE_KAFKA:
		connection, err := connection.MakeKafkaConnection(job.Source.Connection)
		if err != nil {
			stats.ErrCreateKafkaReceiver.Incr()
			return err
		}

		receiver, err := kafka.NewKafkaReceiver(job, *connection, svcMaster, eventsProducer)
		if err != nil {
			stats.ErrCreateKafkaReceiver.Incr()
			return err
		}

		go receiver.ReceiveMessages()
	default:
		stats.ErrCreateReceiver.Incr()
		mlog.Error("Unknown source selected: %v", job.Source.Name)
		os.Exit(-1)
	}

	stats.StartWorker.Incr()
	return nil
}

func getStreamingJob() (*cluster.StartJob, error) {
	job := os.Getenv("DPE_STREAM_WORKER_JOB")
	if job == "" {
		return nil, fmt.Errorf("Please set DPE_STREAM_WORKER_JOB!")
	}

	mlog.Debug("Worker job json: %v", job)

	decoded, err := b64.StdEncoding.DecodeString(job)
	if err != nil {
		return nil, err
	}

	mlog.Debug("Decoded: %v", string(decoded))

	var startJob cluster.StartJob
	err = json.Unmarshal(decoded, &startJob)
	if err != nil {
		mlog.Error("Failed to unmarshal: %v", err)
		return nil, err
	}

	mlog.Debug("Job: %v", startJob)
	return &startJob, nil
}
