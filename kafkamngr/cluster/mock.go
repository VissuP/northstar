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

import "fmt"

type MockCluster struct{}

func (m MockCluster) HasTopic(topicName string) (bool, error) {
	if topicName == "existingTopic" {
		return true, nil
	}
	return false, nil
}

func (m MockCluster) CreateTopic(topicName string, parallelism int,
	replication int) error {
	if topicName == "" {
		return fmt.Errorf("Topic name is empty")
	}
	return nil
}

func (m MockCluster) GetTopicNames() ([]string, error) {
	return []string{"topic1", "topic2", "topic3"}, nil
}

func (m MockCluster) UpdateNumberOfPartitions(serviceName string,
	topicName string, nPartitions int) error {
	return nil
}
