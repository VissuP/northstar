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
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/pkg/msgq"
	"github.com/verizonlabs/northstar/pkg/service_master"
	"github.com/verizonlabs/northstar/dpe-stream/config"
	"github.com/verizonlabs/northstar/dpe-stream/master/cluster"
	"github.com/verizonlabs/northstar/dpe-stream/master/connection"
	"github.com/verizonlabs/northstar/dpe-stream/worker/events"
	"sync/atomic"
	"time"
)

type KafkaReceiver struct {
	topicName      string
	job            *cluster.StartJob
	svcMaster      *service_master.ServiceMaster
	consumer       msgq.MsgQConsumer
	eventsProducer events.EventsProducer
}

func NewKafkaReceiver(job *cluster.StartJob,
	connection connection.KafkaConnection,
	svcMaster *service_master.ServiceMaster,
	eventsProducer events.EventsProducer) (*KafkaReceiver, error) {
	msgQ, err := msgq.NewMsgQ(connection.Topic+"_"+job.AccountId, connection.Brokers, connection.ZK)
	if err != nil {
		mlog.Error("Error to create msgq: %v", err.Error())
		return nil, err
	}

	mlog.Debug("Creating consumer for topic: %v", connection.Topic)
	consumer, err := msgQ.NewConsumer(&msgq.ConsumerConfig{
		TopicName: connection.Topic,
	})
	if err != nil {
		mlog.Error("Error to create msgq consumer: %v", err.Error())
		return nil, err
	}

	for i := 0; i < len(job.Functions); i++ {
		if err = (&(job.Functions[i])).Decode(); err != nil {
			return nil, err
		}
	}

	return &KafkaReceiver{job: job,
		topicName:      connection.Topic,
		svcMaster:      svcMaster,
		consumer:       consumer,
		eventsProducer: eventsProducer}, nil
}

func (r *KafkaReceiver) ReceiveMessages() {
	tickChan := time.NewTicker(time.Duration(config.MsgInterval) * time.Second).C
	var cps uint64 = 0
	for {
		select {
		case event := <-r.consumer.Receive():
			if event.Err != nil {
				mlog.Error(event.Err.Error())
				continue
			}

			worker, err := NewKafkaWorker(r.job, event, r.consumer, r.eventsProducer)
			if err != nil {
				mlog.Error("Failed to create kafka worker: %v", err)
				continue
			}

			if err := r.svcMaster.Dispatch(r.topicName, worker); err != nil {
				mlog.Error(err.Error())
			}
			atomic.AddUint64(&cps, 1)
		case <-tickChan:
			val := atomic.LoadUint64(&cps)
			atomic.SwapUint64(&cps, 0)
			if val == 0 {
				mlog.Alarm("No MSGs on topic %s for %d seconds", r.topicName, config.MsgInterval)
			}
			mlog.Info("CPS on topic: %s is  %.3f",
				r.topicName, float64(val)/float64(config.MsgInterval))
		}
	}
}
