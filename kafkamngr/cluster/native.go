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
	"fmt"
	"github.com/lysu/kazoo-go"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"time"
)

type NativeKafka struct {
	MsgQAdmin *MsgQAdmin
	Client    *kazoo.Kazoo
}

func NewNativeKafka(zkUrl string, zkTimeout int) (KafkaCluster, error) {
	mlog.Info("Native Kafka connection to %s", zkUrl)
	conf := kazoo.NewConfig()
	conf.Timeout = time.Duration(zkTimeout) * time.Millisecond
	kz, err := kazoo.NewKazooFromConnectionString(zkUrl, conf)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Zookeeper: %v", err)
	}

	msgqAdmin, err := NewMsgQAdmin()
	if err != nil {
		return nil, fmt.Errorf("Failed to create msgQ admin: %v", err)
	}

	return &NativeKafka{Client: kz, MsgQAdmin: msgqAdmin}, nil
}

func (n *NativeKafka) HasTopic(topicName string) (bool, error) {
	topic := n.Client.Topic(topicName)
	exists, err := topic.Exists()
	if err != nil {
		return false, fmt.Errorf("Unexpected error (%v) while checking if topic exists", err)
	}
	return exists, nil
}

func (n *NativeKafka) CreateTopic(topicName string, partitions int, replication int) error {
	mlog.Info("Creating topic %s with partitions %d and replication %d", topicName, partitions, replication)
	err := n.Client.CreateTopic(topicName, partitions, replication, make(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("Unexpected error (%v) creating topic %s", err, topicName)
	}
	return nil
}

func (n *NativeKafka) GetTopicNames() ([]string, error) {
	topics, err := n.Client.Topics()
	if err != nil {
		return nil, fmt.Errorf("Unexpected error (%v) while retrieving topics", err)
	}

	names := make([]string, 0)
	for _, topic := range topics {
		names = append(names, topic.Name)
	}

	return names, nil
}

func (n *NativeKafka) UpdateNumberOfPartitions(serviceName string,
	topicName string,
	nPartitions int) error {
	partititons, err := n.MsgQAdmin.Queue.Client.Partitions(topicName)
	totalNumberOfPartitions := int32(len(partititons))
	reqNumberOfPartitions := int32(nPartitions)
	mlog.Debug("Required partitions: %d, total partition count: %d",
		reqNumberOfPartitions, totalNumberOfPartitions)
	if err != nil {
		return err
	}

	var i int32 = 0
	nPartititonsToBL := totalNumberOfPartitions - reqNumberOfPartitions
	if nPartititonsToBL > 0 {
		mlog.Info("BL last %d partitions", nPartititonsToBL)
		for i = reqNumberOfPartitions; i < totalNumberOfPartitions; i++ {
			mlog.Info("Blacklisting partition: %d", i)
			n.MsgQAdmin.markSingleBLPart(i, serviceName, topicName, true)
			n.MsgQAdmin.markSingleBLPart(i, serviceName, topicName, false)
		}
	}

	mlog.Info("WL first %d partitions", reqNumberOfPartitions)
	for i = 0; i < reqNumberOfPartitions; i++ {
		mlog.Info("Whitelisting partition: %d", i)
		n.MsgQAdmin.markSingleWLPart(i, serviceName, topicName, true)
		n.MsgQAdmin.markSingleWLPart(i, serviceName, topicName, false)
	}

	nPartitionsToAdd := reqNumberOfPartitions - totalNumberOfPartitions
	if nPartitionsToAdd > 0 {
		mlog.Info("Adding %d new partitions", nPartitionsToAdd)
		return n.Client.AddPartitions(topicName, nPartitions, "", true)
	}

	return nil
}
