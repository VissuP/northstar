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
)

type ListTopicsCmd struct {
	client *client.KafkaMngrClient
	cmd    *flag.FlagSet
}

func NewListTopics(client *client.KafkaMngrClient) commands.Command {
	cmd := flag.NewFlagSet("topics-list", flag.ExitOnError)
	return &ListTopicsCmd{client: client, cmd: cmd}
}

func (d *ListTopicsCmd) Run(args []string) error {
	d.cmd.Parse(args)

	if !d.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	resp, err := d.client.GetTopics()
	if err != nil {
		return err
	}

	for _, topic := range resp {
		fmt.Println(topic)
	}

	return nil
}
