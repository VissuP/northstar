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

package kafka

import (
	"fmt"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/msgq"
	"github.com/verizonlabs/northstar/dpe-stream/master/cluster"
	"github.com/verizonlabs/northstar/dpe-stream/worker/events"
	"github.com/verizonlabs/northstar/dpe-stream/worker/execution"
)

type KafkaWorker struct {
	execution execution.Execution
	job       *cluster.StartJob
	event     *msgq.ConsumerEvent
	consumer  msgq.MsgQConsumer
}

func NewKafkaWorker(job *cluster.StartJob,
	event *msgq.ConsumerEvent,
	consumer msgq.MsgQConsumer,
	eventsProducer events.EventsProducer) (*KafkaWorker, error) {
	luaExecution, err := execution.NewLuaExecution(eventsProducer)
	if err != nil {
		return nil, err
	}
	return &KafkaWorker{
		execution: luaExecution,
		job:       job,
		event:     event,
		consumer:  consumer,
	}, nil
}

func (s *KafkaWorker) Run(n int) error {
	mlog.Debug("Starting worker %d", n)

	terminate, err := s.execution.ExecuteJob(s.event.Value, s.job)
	if err != nil {
		mlog.Error("Failed to execute functions: %v", err)
		return err
	}

	if terminate {
		return fmt.Errorf("Streaming processing ended, should be shutting down worker")
	}

	mErr := s.consumer.SetAckOffset(s.event.Offset)
	if mErr != nil {
		mlog.Error("Failed to ack offset: %v", mErr)
		return fmt.Errorf(mErr.Error())
	}

	return nil
}
