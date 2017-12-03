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

package events

import (
	"encoding/json"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/msgq"
	"github.com/verizonlabs/northstar/dpe-stream/master/cluster"
	"github.com/verizonlabs/northstar/dpe-stream/worker/stats"
)

const (
	DPE_STREAM_SERVICE_NAME = "dpe-stream"
	DPE_STREAM_OUTPUT_TOPIC = "dpe-stream-output"
)

type KafkaEventsProducer struct {
	msgQ          msgq.MessageQueue
	kafkaProducer msgq.MsgQProducer
}

func NewKafkaEventsProducer() (*KafkaEventsProducer, error) {
	msgQ, err := msgq.NewMsgQ(DPE_STREAM_SERVICE_NAME, nil, nil)
	if err != nil {
		return nil, err
	}

	kafkaProducer, err := msgQ.NewProducer(&msgq.ProducerConfig{
		TopicName:   DPE_STREAM_OUTPUT_TOPIC,
		Partitioner: msgq.RoundRobinPartitioner,
	})
	if err != nil {
		return nil, err
	}

	return &KafkaEventsProducer{kafkaProducer: kafkaProducer}, nil
}

func (p *KafkaEventsProducer) StreamOutput(job *cluster.StartJob,
	stdout string,
	stderr string,
	result string) error {
	event := &StreamOutputEvent{AccountId: job.AccountId,
		JobId:        job.JobId,
		InvocationId: job.InvocationId,
		StdOut:       stdout,
		StdErr:       stderr,
		Result:       result}

	outputByte, err := json.Marshal(event)
	if err != nil {
		stats.ErrStreamOutput.Incr()
		return err
	}

	mlog.Debug("Sending to Kafka topic %v: %v", DPE_STREAM_OUTPUT_TOPIC, string(outputByte))
	err = p.kafkaProducer.Send(outputByte)
	if err != nil {
		stats.ErrStreamOutput.Incr()
		return err
	}

	stats.StreamOutput.Incr()
	return nil
}
