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

package connection

import (
	"fmt"
	"github.com/verizonlabs/northstar/pkg/msgq"
)

type KafkaConnection struct {
	Topic   string   `json:"topic,omitempty"`
	Brokers []string `json:"brokers,omitempty"`
	ZK      []string `json:"zk,omitempty"`
}

func NewKafkaConnection() *KafkaConnection {
	return &KafkaConnection{}
}

func (c *KafkaConnection) Validate() error {
	if c.Topic == "" {
		return fmt.Errorf("Topic name is empty")
	}

	if len(c.Brokers) < 1 {
		return fmt.Errorf("Brokers count is is less than one")
	}

	if len(c.ZK) < 1 {
		return fmt.Errorf("ZK count is less than one")
	}

	return nil
}

func (k *KafkaConnection) GetNumberOfWorkers() (int, error) {
	msgQ, err := msgq.NewMsgQ("", k.Brokers, k.ZK)
	if err != nil {
		return 0, err
	}

	nPartitions, err := msgQ.Partitions(k.Topic)
	if err != nil {
		return 0, err
	}

	return nPartitions, nil
}
