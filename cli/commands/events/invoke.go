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

	"encoding/json"
	"errors"
	"fmt"
	"github.com/verizonlabs/northstar/cli/commands"
	"github.com/verizonlabs/northstar/cli/util"
	"github.com/verizonlabs/northstar/processing/events/client"
	"github.com/verizonlabs/northstar/processing/events/model"
)

type InvokeEventCmd struct {
	client *client.EventsClient
	cmd    *flag.FlagSet
	id     *string
	args   *string
}

func UnmarshalArgs(args string) (map[string]interface{}, error) {
	kv := make(map[string]interface{})
	err := json.Unmarshal([]byte(args), &kv)
	return kv, err
}

func NewInvokeEvent(client *client.EventsClient) commands.Command {
	cmd := flag.NewFlagSet("events-invoke", flag.ExitOnError)
	id := cmd.String("id", "", "The event id")
	args := cmd.String("args", "{}", "The event arguments")

	return &InvokeEventCmd{client: client,
		cmd:  cmd,
		id:   id,
		args: args}
}

func (invoke *InvokeEventCmd) Run(args []string) error {
	invoke.cmd.Parse(args)

	if !invoke.cmd.Parsed() {
		return errors.New("Failed to parse cmd")
	}

	if *invoke.id == "" {
		return errors.New("Please set an id using -id.")
	}

	arguments, err := UnmarshalArgs(*invoke.args)
	if err != nil {
		return err
	}

	options := &model.Options{Args: arguments}
	out, mErr := invoke.client.InvokeEvent(util.GetAccountID(), *invoke.id, options)
	if mErr != nil {
		return mErr
	}

	fmt.Println(out)
	return nil
}
