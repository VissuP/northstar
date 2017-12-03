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
	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/dpe-stream/master/model"
	"strings"
)

func GetNumberOfWorkers(src *model.Source) (int, error) {
	switch src.Name {
	case model.SOURCE_KAFKA:
		connection, err := MakeKafkaConnection(src.Connection)
		if err != nil {
			return 0, err
		}

		return connection.GetNumberOfWorkers()
	default:
		errM := fmt.Sprintf("Unknown source selected: %v", src.Name)
		return 0, fmt.Errorf(errM)
	}
}

func MakeKafkaConnection(connection interface{}) (*KafkaConnection, error) {
	kafka := NewKafkaConnection()
	conn, ok := connection.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid kafka connection description")
	}

	topic, ok := conn["Topic"].(string)
	if !ok {
		return nil, errors.New("invalid kafka topic description")
	}
	kafka.Topic = topic

	brokers, ok := conn["Brokers"].(string)
	if !ok {
		return nil, errors.New("invalid kafka brokers description")
	}
	for _, broker := range strings.Split(brokers, ",") {
		kafka.Brokers = append(kafka.Brokers, strings.TrimSpace(broker))
	}

	zk, ok := conn["ZK"].(string)
	if !ok {
		return nil, errors.New("invalid zookeeper nodes description")
	}
	for _, node := range strings.Split(zk, ",") {
		kafka.ZK = append(kafka.ZK, strings.TrimSpace(node))
	}

	return kafka, nil
}
