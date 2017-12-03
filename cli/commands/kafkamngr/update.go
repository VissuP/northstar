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

package kafkamngr

import (
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/kafkamngr/client"
	"github.com/verizonlabs/northstar/kafkamngr/model"
)

type UpdateTopicsCmd struct {
	client      *client.KafkaMngrClient
	cmd         *flag.FlagSet
	service     *string
	topic       *string
	partitions  *int
	replication *int
}

func NewUpdateTopics(client *client.KafkaMngrClient) commands.Command {
	cmd := flag.NewFlagSet("topics-update", flag.ExitOnError)
	service := cmd.String("service", "test", "The service name")
	topic := cmd.String("topic", "test", "The topic name")
	partitions := cmd.Int("partitions", 1, "The number of partitions")
	replication := cmd.Int("replication", 3, "The replication factor")
	return &UpdateTopicsCmd{client: client,
		cmd:         cmd,
		service:     service,
		topic:       topic,
		partitions:  partitions,
		replication: replication}
}

func (output *UpdateTopicsCmd) Run(args []string) error {
	output.cmd.Parse(args)

	if !output.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	topic := &model.Topic{Name: *output.topic,
		Partitions:  *output.partitions,
		Replication: *output.replication}

	err := output.client.UpdateTopic(*output.service, *output.topic, topic)
	if err != nil {
		return err
	}

	fmt.Println("Topic updated")
	return nil
}
