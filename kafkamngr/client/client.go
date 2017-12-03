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

package client

import (
	"encoding/json"
	"fmt"
	lb "github.com/verizonlabs/northstar/pkg/httpclientlb"
	"github.com/verizonlabs/northstar/pkg/management"
	"github.com/verizonlabs/northstar/pkg/mlog"
	"github.com/verizonlabs/northstar/kafkamngr/model"
	"github.com/verizonlabs/northstar/kafkamngr/util"
)

const BASE_URI = util.KafkaMngrBasePath + "/topics"

type KafkaMngrClient struct {
	lbClient *lb.LbClient
}

func NewKafkaMngrClient() (*KafkaMngrClient, error) {
	url, err := util.GetKafkaMngrBaseUrl()
	if err != nil {
		mlog.Error("Failed to get kafka mngr base url with error: %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	lbClient, err := lb.GetClient(url)
	if err != nil {
		mlog.Info("Failed to create kafkamngr client with error: %s", err.Error())
		return nil, management.GetInternalError(err.Error())
	}

	return &KafkaMngrClient{lbClient: lbClient}, nil
}

func (client *KafkaMngrClient) CreateTopic(topic *model.Topic) *management.Error {
	path := fmt.Sprintf("%s", BASE_URI)
	_, err := client.lbClient.PostJSON(path, topic)
	if err != nil {
		mlog.Error("Kafka Mngr client: Error creating topic by name: %s", err.Error())
		return err
	}
	return nil
}

func (client *KafkaMngrClient) GetTopics() ([]string, *management.Error) {
	path := fmt.Sprintf("%s", BASE_URI)
	data, err := client.lbClient.Get(path)
	if err != nil {
		return nil, err
	}

	var out []string
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, management.GetInternalError(err.Error())
	}

	return out, nil
}

func (client *KafkaMngrClient) UpdateTopic(serviceName string,
	topicName string,
	topic *model.Topic) *management.Error {
	path := fmt.Sprintf("%s/%s/%s", BASE_URI, serviceName, topicName)
	_, err := client.lbClient.PutJSON(path, topic)
	if err != nil {
		return err
	}
	return nil
}
