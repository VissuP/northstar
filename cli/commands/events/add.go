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
	"flag"

	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/data/events/client"
	"github.com/verizonlabs/northstar/data/events/model"
)

type AddEventCmd struct {
	client *client.EventsClient
	cmd    *flag.FlagSet
	id     *string
	name   *string
}

func NewAddEvent(client *client.EventsClient) commands.Command {
	cmd := flag.NewFlagSet("events-add", flag.ExitOnError)
	id := cmd.String("id", "", "The event id")
	name := cmd.String("name", "", "The event name")

	return &AddEventCmd{client: client,
		id:   id,
		cmd:  cmd,
		name: name}
}

func (output *AddEventCmd) Run(args []string) error {
	output.cmd.Parse(args)

	if !output.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *output.name == "" {
		return errors.New("Please set a name using -name.")
	}

	data := &model.EventData{Id: *output.id, Name: *output.name}
	id, err := output.client.AddEvent(util.GetAccountID(), data)
	if err != nil {
		return err
	}

	fmt.Printf("Event %s added", id)
	return nil
}
